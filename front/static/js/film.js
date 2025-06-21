async function loadProdData(login) {
    try {
        const response = await fetch('/admin/prod_data', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ login })
        });
        
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }
        
        return await response.json();
    } catch (error) {
        console.error('User data load error:', error);
        return null;
    }
}

async function loadRegisData(login) {
    try {
        const response = await fetch('/admin/regis_data', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ login })
        });
        
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }
        
        return await response.json();
    } catch (error) {
        console.error('User data load error:', error);
        return null;
    }
}

async function addProdus(logg) {
    globalThis.logP = logg
}

async function addRegisse(logg) {
    globalThis.logR = logg
}

document.getElementById("searchProdBtn").addEventListener("click", async function() {
    const login1 = document.getElementById("searchProducer").value;
    const userData1 = await loadProdData(login1);
    if (userData1) {
        const resultsContainer1 = document.getElementById("prodResults");
        if (resultsContainer1) {
            resultsContainer1.innerHTML = `
                <div class="card mt-3">
                    <div class="card-body">
                        <h5 class="card-title">${userData1.login}</h5>
                        <p class="card-text">
                            <p><strong>ФИО:</strong> ${userData1.last_name} ${userData1.first_name} ${userData1.middle_name}</p>
                            <strong>Дата рождения:</strong> ${userData1.birth_date}<br>
                            <strong>Роль:</strong> ${userData1.role_name}<br>
                            <strong>Опыт работы:</strong> ${userData1.experience_years}
                        </p>
                        <button id="addProdus" class="btn btn-danger" data-login="${userData1.login}">
                            <i class="fas fa-trash"></i> Прикрепить продюсера к фильму
                        </button>
                    </div>
                </div>
            `;
            document.getElementById('addProdus').addEventListener('click', addProdus(login1));
        }
    } else {
        console.error("Не удалось загрузить данные пользователя");
        alert("Пользователь не найден или произошла ошибка");
    }
});

document.getElementById("searchRegisBtn").addEventListener("click", async function() {
    const login2 = document.getElementById("searchRegisser").value;
    const userData2 = await loadRegisData(login2);
    if (userData2) {
        const resultsContainer2 = document.getElementById("regisResults");
        if (resultsContainer2) {
            resultsContainer2.innerHTML = `
                <div class="card mt-3">
                    <div class="card-body">
                        <h5 class="card-title">${userData2.login}</h5>
                        <p class="card-text">
                            <p><strong>ФИО:</strong> ${userData2.last_name} ${userData2.first_name} ${userData2.middle_name}</p>
                            <strong>Дата рождения:</strong> ${userData2.birth_date}<br>
                            <strong>Роль:</strong> ${userData2.role_name}<br>
                            <strong>Опыт работы:</strong> ${userData2.experience_years}
                        </p>
                        <button id="addRegisse" class="btn btn-danger" data-login="${userData2.login}">
                            <i class="fas fa-trash"></i> Прикрепить режиссера к фильму
                        </button>
                    </div>
                </div>
            `;
            document.getElementById('addRegisse').addEventListener('click', addRegisse(login2));
        }
    } else {
        console.error("Не удалось загрузить данные пользователя");
        alert("Пользователь не найден или произошла ошибка");
    }
});


document.getElementById("addFilm").addEventListener('click', async function() {
    event.preventDefault();
    
    const movieData = {
        title: document.getElementById('movieTitle').value,
        producer: logP,
        director: logR
    };

    try {
        const response = await fetch('/admin/create_film', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(movieData)
        });

        const result = await response.json();
        
        if (response.ok) {
            alert('Фильм успешно добавлен!');
            this.reset();
            document.getElementById('prodResults').innerHTML = '';
            document.getElementById('regisResults').innerHTML = '';
        } else {
            throw new Error(result.error || 'Ошибка сервера');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка при добавлении фильма: ' + error.message);
    }
})