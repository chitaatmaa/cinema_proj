        document.getElementById('regUser').addEventListener('click', function() {
            registerUser();
        });

        function registerUser() {
            const roleId = document.getElementsByName('role_id')[0].value;
            const firstName = document.getElementsByName('first_name')[0].value;
            const lastName = document.getElementsByName('last_name')[0].value;
            const middleName = document.getElementsByName('middle_name')[0].value;
            const birthDate = document.getElementsByName('birth_date')[0].value;
            const photo = document.getElementsByName('photo')[0].files[0];
            const experienceYears = document.getElementsByName('experience_years')[0].value;
            const login = document.getElementsByName('login')[0].value;
            const pass = document.getElementsByName('pass')[0].value;
            
            if (!firstName || !lastName || !birthDate || !login || !pass) {
                showMessage('Пожалуйста, заполните все обязательные поля', 'error');
                return;
            }

            const formData = new FormData();
            formData.append('role_id', roleId);
            formData.append('first_name', firstName);
            formData.append('last_name', lastName);
            formData.append('middle_name', middleName);
            formData.append('birth_date', birthDate);
            formData.append('experience_years', experienceYears);
            formData.append('login', login);
            formData.append('pass', pass);
            
            if (photo) {
                formData.append('photo', photo);
            }

            fetch('/register', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    window.location.href = '/auth';
                } else {
                    return response.text().then(text => { throw new Error(text) });
                }
            })
            .catch(error => {
                showMessage('Ошибка регистрации: ' + error.message, 'error');
                console.error('Ошибка:', error);
            })
            .finally(() => {
                button.disabled = false;
                button.innerHTML = 'Зарегистрироваться';
            });
        }

        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.className = type;
        }