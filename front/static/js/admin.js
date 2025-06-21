document.addEventListener('DOMContentLoaded', () => {    
        const login = sessionStorage.getItem('UL');
        if (login) {
            document.getElementById("admin_login").textContent=login;
            const imgX = document.getElementById('userPhoto');
            const img = new Image();
            img.src = `/admin/user_photo?login=${encodeURIComponent(login)}`;
            img.onload = () => {
                imgX.src = `/admin/user_photo?login=${encodeURIComponent(login)}`
            };
            img.onerror = () => {
                imgX.src = "../static/images/default_user_photo.jpeg";
            };
            sessionStorage.removeItem('UL');
        } else { window.location.href = '/auth'; }
        const dropdown = document.getElementById('adminDropdown');
        dropdown.addEventListener('click', function(e) {
            const menu = this.querySelector('.dropdown-menu');
            menu.classList.toggle('show');
            e.stopPropagation();
        });
        document.addEventListener('click', function(e) {
            if (!dropdown.contains(e.target)) {
                const menu = dropdown.querySelector('.dropdown-menu');
                menu.classList.remove('show');
            }
        });
        document.getElementById('userForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const userRole = document.getElementsByName('role_id')[0].value;
            const password = document.getElementById('password').value;
            const login = document.getElementById('login').value;
            const birthDate = document.getElementById('birthDate').value;
            const middleName = document.getElementById('middleName').value;
            const lastName = document.getElementById('lastName').value;
            const firstName = document.getElementById('firstName').value;
            const experience_years = document.getElementsByName('experience_years')[0].value;
            const photo = document.getElementsByName('photo')[0].files[0];
            
            const formData = new FormData();
            formData.append('role_id', userRole);
            formData.append('first_name', firstName);
            formData.append('last_name', lastName);
            formData.append('middle_name', middleName);
            formData.append('birth_date', birthDate);
            formData.append('experience_years', experience_years);
            formData.append('login', login);
            formData.append('pass', password);
            if (photo) {
                formData.append('photo', photo);
            }


            fetch('/admin', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    alert("Успешно зарегистрирован пользователь ");
                } else {
                    return response.text().then(text => { throw new Error(text) });
                }
            })
            .catch(error => {
                showMessage('Ошибка регистрации: ' + error.message, 'error');
                console.error('Ошибка:', error);
            })
        });
        
})

async function loadUserData(login) {
    try {
        const response = await fetch('/admin/user_data', {
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

document.getElementById("searchUserBtn").addEventListener("click", async function() {
    const login = document.getElementById("searchLogin").value;
    const userData = await loadUserData(login);
    if (userData) {
        console.log(userData)
        const resultsContainer = document.getElementById("userResults");
        if (resultsContainer) {
            resultsContainer.innerHTML = `
                <div class="card mt-3">
                    <div class="card-body">
                        <h5 class="card-title">${userData.login}</h5>
                        <p class="card-text">
                            <p><strong>ФИО:</strong> ${userData.last_name} ${userData.first_name} ${userData.middle_name}</p>
                            <strong>Дата рождения:</strong> ${userData.birth_date}<br>
                            <strong>Роль:</strong> ${userData.role_name}
                        </p>
                        <button id="deleteUserBtn" class="btn btn-danger" data-login="${userData.login}">
                            <i class="fas fa-trash"></i> Удалить пользователя
                        </button>
                    </div>
                </div>
            `;
            document.getElementById('deleteUserBtn').addEventListener('click', deleteUser);
        }
    } else {
        console.error("Не удалось загрузить данные пользователя");
        alert("Пользователь не найден или произошла ошибка");
    }
});

async function deleteUser() {
    const login = this.getAttribute('data-login');
    const resultsContainer = document.getElementById('userSearchResult');
    if (!confirm(`Вы уверены, что хотите удалить пользователя ${login}?`)) {
        return;
    }

    try {
        const response = await fetch('/admin/delete', {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ login })
        });

        const data = await response.json();

        if (!response.ok) {
            resultsContainer.innerHTML = `<div class="alert alert-danger">${data.error || 'Ошибка удаления'}</div>`;
            return;
        } else {
            const results = document.getElementById("userResults");
            document.getElementById("searchLogin").value='';
            results.remove(); 
        }
    } catch (error) {
        console.log(error)
    }
}