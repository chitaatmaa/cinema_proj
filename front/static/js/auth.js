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
                    window.location.href = '/producer';                
                } else {
                    window.location.href = '/regisser';                
                }
            } else {
                alert(data.error || 'Auth failed');
            }
        } catch (err) {
            console.error('Error:', err);
        }
    });
});