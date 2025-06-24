<<<<<<< HEAD
document.addEventListener('DOMContentLoaded', () => {
    const ab1 = document.getElementById("authb1");
    const ab2 = document.getElementById("authb2");

    ab2.addEventListener('click', function () {
        window.location.replace('/register')
    })

    ab1.addEventListener('click', async (e) => {
        e.preventDefault();

        const logine = document.getElementById('login').value;
        const passe = document.getElementById('password').value;
        
        const formData = {
            login: logine,
            password: passe
        };

        try {
            const response = await fetch('/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData)
            });
            
            const data = await response.json();
            
            if (response.ok) {
                if (data.role_id == '0') {
                    sessionStorage.setItem('UL', logine);
                    window.location.href = '/admin';
                } else if (data.role_id == '1') {
                    sessionStorage.setItem('UL', logine);
                    window.location.href = `/producer?login=${encodeURIComponent(logine)}`;              
                } else {
                    sessionStorage.setItem('UL', logine);
                    window.location.href = `/regisser?login=${encodeURIComponent(logine)}`;              
                }
            } else {
                alert(data.error || 'Auth failed');
            }
        } catch (err) {
            console.error('Error:', err);
        }
    });
=======
document.addEventListener('DOMContentLoaded', () => {
    const ab1 = document.getElementById("authb1");
    const ab2 = document.getElementById("authb2");

    ab2.addEventListener('click', function () {
        window.location.replace('/register')
    })

    ab1.addEventListener('click', async (e) => {
        e.preventDefault();

        const logine = document.getElementById('login').value;
        const passe = document.getElementById('password').value;
        
        const formData = {
            login: logine,
            password: passe
        };

        try {
            const response = await fetch('/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData)
            });
            
            const data = await response.json();
            
            if (response.ok) {
                if (data.role_id == '0') {
                    sessionStorage.setItem('UL', logine);
                    window.location.href = '/admin';
                } else if (data.role_id == '1') {
                    sessionStorage.setItem('UL', logine);
                    window.location.href = `/producer?login=${encodeURIComponent(logine)}`;              
                } else {
                    sessionStorage.setItem('UL', logine);
                    window.location.href = `/regisser?login=${encodeURIComponent(logine)}`;              
                }
            } else {
                alert(data.error || 'Auth failed');
            }
        } catch (err) {
            console.error('Error:', err);
        }
    });
>>>>>>> 741fb8c1d90e4e1b14d660659e9dfa19713f6128
});