document.addEventListener('DOMContentLoaded', function() {
    const login = sessionStorage.getItem('UL');
    if (!login) {window.location.href='/auth';}
    sessionStorage.removeItem('UL');

    document.getElementById("regis_login").textContent=login;
    
    // Данные для текущего фильма
    let currentFilm = null;
    let filmGroups = [];
    let actors =[];

    let filmActors = [];
    const collapsibles = document.querySelectorAll('.collapsible');

    // Фильмы
    const filmSelect = document.getElementById('film-select');
    const selectFilmBtn = document.getElementById('select-film-btn');
    const filmTitle = document.getElementById('film-title');
    const filmInfo = document.getElementById('film-info');

    async function startFilm() {
        if (!currentFilm) {
            alert('Сначала выберите фильм');
            return;
        }
        
        if (filmGroups.length === 0 && actors.length === 0) {
            alert('Добавьте хотя бы одну группу или актера');
            return;
        }

        const requestData = {
            film_id: currentFilm.id,
            groups: filmGroups.map(g => ({
                group_id: g.id,
                cost: g.cost
            })),
            actors: actors.map(a => ({
                actor_id: a.id,
                cost: a.cost1,
                character_name: a.scenic
            }))
        };
        console.log(requestData)

        try {
            const response = await fetch('/regisser/start_film', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(requestData)
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Ошибка сервера');
            }

            const result = await response.json();
            showNotification('summary-notification', result.message);
            
            currentFilm = null;
            filmGroups = [];
            actors = [];
            filmTitle.textContent = '';
            filmInfo.classList.add('hidden');
            updateSummary();
            updateSummary1();
            updateTotalCost();
            
        } catch (error) {
            console.error('Ошибка запуска фильма:', error);
            alert('Ошибка: ' + error.message);
        }
    }

    async function saveGroup() {
        const groupName = document.getElementById('group-name').value;
        const groupSize = parseInt(document.getElementById('group-size').value);
        
        if (!groupName || groupSize <= 0) {
            alert('Please fill all fields correctly');
            return;
        }

        try {
            const response = await fetch('/regisser/add_group', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: groupName,
                    count: groupSize
                })
            });

            const responseData = await response.json();
            if (!response.ok) {
                throw new Error(responseData.error || 'Server error');
            }
            
            updateGroupDropdown(result.group);
            window.location.reload();
            resetGroupForm();
        } catch (error) {
            console.error('Error saving group:', error);
            alert('Error: ' + error.message);
        }
    }

    async function saveActor() {
        const logA = document.getElementById('actor-login').value;
        const firstA = document.getElementById('actor-first-name').value;
        const lastA = document.getElementById('actor-last-name').value;
        const middleA = document.getElementById('actor-middle-name').value;
        const birthA = document.getElementById('actor-birthdate').value;
        const expA = parseInt(document.getElementById('actor-experience').value);
        const mailA = document.getElementById('actor-email').value;
        const phoneA = document.getElementById('actor-phone').value;

        
        if (!logA || !firstA || !lastA || !middleA || !birthA || !expA || !mailA || !phoneA) {
            console.log(logA, firstA, lastA, middleA, birthA, expA, mailA, phoneA)
            alert('Please fill all fields correctly');
            return;
        }
        console.log(logA, firstA, lastA, middleA, birthA, expA, mailA, phoneA)

        try {
            const response = await fetch('/regisser/add_actor', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    login: logA,
                    first_name: firstA,
                    last_name: lastA,
                    middle_name: middleA,
                    birth_date: birthA,
                    experience: expA,
                    email: mailA,
                    phone: phoneA,
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || 'Server error');
            }

            const errorData = await response.json();
            window.location.reload();
            resetActorForm();
        } catch (error) {
            console.error('Error saving group:', error);
            alert('Error: ' + error.message);
        }
    }

    function updateGroupDropdown(group) {
        const select = document.getElementById('existing-group');
        const option = document.createElement('option');
        option.value = group.id;
        option.textContent = `${group.name} (${group.count} members)`;
        select.appendChild(option);
    }

    // Группы
    const saveGroupBtn = document.getElementById('save-group-btn');
    let groupAction = document.getElementById('group-action');
    const newGroupSection = document.getElementById('new-group-section');
    const existingGroupSection = document.getElementById('existing-group-section');
    const groupsTableBody = document.getElementById('groups-table-body');

    saveGroupBtn.addEventListener('click', function() {
        if (groupAction.value === "add") {
            saveGroup();
        } else {
            addExistingGroupToSummary();
        }
    });

   function addExistingGroupToSummary() {
        const groupSelect = document.getElementById('existing-group');
        const selectedOption = groupSelect.options[groupSelect.selectedIndex];
        const groupId = groupSelect.value;
        const costInput = document.getElementById('group-cost');
        
        if (!groupId) {
            alert('Пожалуйста, выберите группу');
            return;
        }
        
        const cost = parseInt(costInput.value) || 0;
        if (cost <= 0) {
            alert('Пожалуйста, укажите стоимость');
            return;
        }
        // Добавляем группу в итоговую информацию
        let group = {
            id: parseInt(groupId),
            name: selectedOption.textContent,
            cost: cost,
        };
        
        // Добавляем группу к текущему фильму
        if (currentFilm) {
            filmGroups.push(group);
            updateSummary();
            updateTotalCost();
            showNotification('group-notification', 'Группа добавлена в фильм');
        } else {
            alert('Пожалуйста, сначала выберите фильм');
        }
        
        // Сбрасываем форму
        groupSelect.selectedIndex = 0;
        costInput.value = '';
    }

    groupAction.addEventListener('change', function() {
        if (this.value === 'add') {
            newGroupSection.style.display = 'block';
            existingGroupSection.classList.add('hidden');
            existingGroupSection.style.display = 'none';
        } else {
            newGroupSection.classList.add('hidden');
            newGroupSection.style.display = 'none';
            existingGroupSection.style.display = 'block';
        }
    });

    // Актеры
    const actorAction = document.getElementById('actor-action');
    const newActorSection = document.getElementById('new-actor-section');
    const existingActorSection = document.getElementById('existing-actor-section');
    const saveActorBtn = document.getElementById('save-actor-btn');
    const actorsTableBody = document.getElementById('actors-table-body');

    const editSummaryBtn = document.getElementById('edit-summary-btn');
    
    collapsibles.forEach(collapsible => {
        collapsible.addEventListener('click', function() {
            this.classList.toggle('active');
            const content = this.nextElementSibling;
            content.classList.toggle('active');
        });
    });
    
    actorAction.addEventListener('change', function() {
        if (this.value === 'add') {
            newActorSection.style.display = 'block';
            newActorSection.classList.remove('hidden');
            existingActorSection.classList.add('hidden');
            existingActorSection.style.display = 'none';
        } else {
            newActorSection.style.display = 'none';
            newActorSection.classList.add('hidden');
            existingActorSection.classList.remove('hidden');
            existingActorSection.style.display = 'block';
        }
    });

    saveActorBtn.addEventListener('click', function() {
        if (actorAction.value === "add") {
            saveActor();
        } else {
            addExistingActorToSummary();
        }
    });
    
    selectFilmBtn.addEventListener('click', function() {
        const filmId = filmSelect.value;
        const selectedOption = filmSelect.options[filmSelect.selectedIndex];
        
        if (filmId) {
            currentFilm = {
                id: parseInt(filmId),
                title: selectedOption.textContent
            };
            filmTitle.textContent = currentFilm.title;
            filmInfo.classList.remove('hidden');
            document.getElementById('summary-notification').style.display = 'none';
            showNotification('film-notification', 'Фильм успешно выбран!');
            updateSummary();
        } else {
            alert('Пожалуйста, выберите фильм');
        }
    });

    function addExistingActorToSummary() {
        const actorSelect = document.getElementById('existing-actor');
        const selectedOption1 = actorSelect.options[actorSelect.selectedIndex];
        const actorId = actorSelect.value;
        const costInput1 = document.getElementById('actor-cost');
        const scenic_name = document.getElementById('character-name').value.toString();
        
        if (!actorId) {
            alert('Пожалуйста, выберите актера');
            return;
        }
        
        const cost1 = parseInt(costInput1.value) || 0;
        if (cost1 <= 0) {
            alert('Пожалуйста, укажите стоимость');
            return;
        }

        if (!scenic_name) {
            alert('Пожалуйста, укажите сценическое имя');
            return;
        }

        let actor = {
            id: parseInt(actorId),
            login: selectedOption1.textContent,
            cost1: cost1,
            scenic: scenic_name,
        };
        
        // Добавляем группу к текущему фильму
        if (currentFilm) {
            actors.push(actor);
            updateSummary1();
            updateTotalCost();
            showNotification('actor-notification', 'Группа добавлена в фильм');
        } else {
            alert('Пожалуйста, сначала выберите фильм');
        }
        
        // Сбрасываем форму
        actorSelect.selectedIndex = 0;
        costInput1.value = '';
    }
    
    editSummaryBtn.addEventListener('click', function() {
        // Открываем первую секцию для редактирования
        const firstCollapsible = document.querySelector('.collapsible');
        firstCollapsible.classList.add('active');
        firstCollapsible.nextElementSibling.classList.add('active');
        
        // Прокручиваем к началу страницы
        window.scrollTo({ top: 0, behavior: 'smooth' });
    });
    
    // Вспомогательные функции
    function showNotification(id, message) {
        const notification = document.getElementById(id);
        notification.textContent = message;
        notification.style.display = 'block';
        
        setTimeout(() => {
            notification.style.display = 'none';
        }, 3000);
    }
    
    function resetGroupForm() {
        document.getElementById('group-name').value = '';
        document.getElementById('group-size').value = '5';
        document.getElementById('group-cost').value = '25000';
        document.getElementById('existing-group').selectedIndex = 0;
        groupAction.value = 'add';
        newGroupSection.style.display = 'block';
        existingGroupSection.classList.add('hidden');
    }
    
    function resetActorForm() {
        document.getElementById('actor-login').value = '';
        document.getElementById('actor-first-name').value = '';
        document.getElementById('actor-last-name').value = '';
        document.getElementById('actor-middle-name').value = '';
        document.getElementById('actor-birthdate').value = '';
        document.getElementById('actor-experience').value = '5';
        document.getElementById('actor-email').value = '';
        document.getElementById('actor-phone').value = '';
        document.getElementById('character-name').value = '';
        document.getElementById('actor-cost').value = '15000';
        document.getElementById('existing-actor').selectedIndex = 0;
        actorAction.value = 'add';
        newActorSection.style.display = 'block';
        existingActorSection.classList.add('hidden');
    }

    document.getElementById('start-film-btn').addEventListener('click', startFilm);
    
    function updateSummary() {
        groupsTableBody.innerHTML = '';
        filmGroups.forEach(group => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${group.name}</td>
                <td>${group.cost.toLocaleString()} руб.</td>
            `;
            groupsTableBody.appendChild(row);
        });
    }

    function updateSummary1() {
        actorsTableBody.innerHTML = '';
        actors.forEach(actor => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${actor.login}</td>
                <td>${actor.cost1.toLocaleString()} руб.</td>
                <td>${actor.scenic}</td>
            `;
            actorsTableBody.appendChild(row);
        });
    }    

    function updateTotalCost() {
        const groupsCost = filmGroups.reduce((sum, group) => sum + group.cost, 0);
        const actorsCost = actors.reduce((sum, actor) => sum + actor.cost1, 0);
        const totalCost = groupsCost + actorsCost;
        document.getElementById('total-cost').textContent = totalCost;
    }
});