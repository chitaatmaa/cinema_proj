const producerLogin = sessionStorage.getItem('UL');
sessionStorage.removeItem('UL');
document.addEventListener('DOMContentLoaded', () => {
    if (!producerLogin) {window.location.replace('/auth');}
    // Загрузка фильмов продюсера
    fetch(`/producer/movies?login=${producerLogin}`)
        .then(response => response.json())
        .then(movies => {
            const select = document.getElementById('movies-select');
            movies.forEach(movie => {
                const option = document.createElement('option');
                option.value = movie.id;
                option.textContent = `${movie.title} (${movie.status})`;
                select.appendChild(option);
            });
            
            select.addEventListener('change', function() {
                if (this.value) {
                    const movie = movies.find(m => m.id == this.value);
                    document.getElementById('budget-section').style.display = 'block';
                    document.getElementById('movie-info').innerHTML = `
                        <p>Текущий бюджет: ${movie.budget || 'не установлен'}</p>
                    `;
                } else {
                    document.getElementById('budget-section').style.display = 'none';
                }
            });
        });
});

function updateBudget() {
    const movieId = parseInt(document.getElementById('movies-select').value);
    const budgetInput = document.getElementById('budget-input');
    const budget = parseInt(budgetInput.value)

    fetch(`/producer/movie/budget`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
            movie_id: movieId, 
            budget: budget 
        })
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById('budget-message').textContent = data.message;
        document.getElementById('budget-message').style.color = 'green';
    })
    .catch(error => {
        document.getElementById('budget-message').textContent = 'Ошибка: ' + error;
        document.getElementById('budget-message').style.color = 'red';
    });
}

function downloadGroupsReport() {
    window.location.href = `/producer/report/groups?login=${producerLogin}`;
}

function downloadDetailedReport() {
    window.location.href = `/producer/report/detailed?login=${producerLogin}`;
}