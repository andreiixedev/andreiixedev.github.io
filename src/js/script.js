// --- Navbar mobil ---
const btn = document.getElementById('toggleMenu');
const mobileNavbar = document.getElementById('mobileNavbar');
const menuLinks = document.querySelectorAll('.menu-links a');
const themeBtn = document.getElementById('themeToggle');
const themeBtnMobile = document.getElementById('themeToggleMobile');
const body = document.body;

function toggleMenu() {
    btn.classList.toggle('open');
    mobileNavbar.classList.toggle('open');
}

btn.addEventListener('click', toggleMenu);

menuLinks.forEach(link => {
    link.addEventListener('click', () => {
        if (mobileNavbar.classList.contains('open')) {
            toggleMenu();
        }
    });
});

// --- Dark/Light Mode ---
function setTheme(light) {
    if (light) {
        body.classList.add('light-mode');
        themeBtn.textContent = "☀️";
        themeBtnMobile.textContent = "☀️";
        localStorage.setItem("theme", "light");
    } else {
        body.classList.remove('light-mode');
        themeBtn.textContent = "🌙";
        themeBtnMobile.textContent = "🌙";
        localStorage.setItem("theme", "dark");
    }
}

function toggleTheme() {
    setTheme(!body.classList.contains('light-mode'));
}

themeBtn.addEventListener('click', toggleTheme);
themeBtnMobile.addEventListener('click', toggleTheme);

// Load theme from localStorage
if (localStorage.getItem("theme") === "light") {
    setTheme(true);
}

if (document.querySelector(".auto-input")) {
    const texts = ["Andreiixe", "Macintosh", "Pikachu", ":)", "hell"];
    const span = document.querySelector(".auto-input");
    let index = 0;

    function showNextText() {
        // Fade out anterior
        span.style.opacity = 0;
        span.style.transform = "translateY(20px) scale(0.8)";

        setTimeout(() => {
            // Schimbă textul
            span.textContent = texts[index];

            // Fade + slide + scale in
            span.style.opacity = 1;
            span.style.transform = "translateY(0) scale(1)";

            // Next text
            index = (index + 1) % texts.length;
            setTimeout(showNextText, 2500); // timp afișare text
        }, 400); // timp animație out
    }

    showNextText();
}


// --- Fetch ultimele commits GitHub ---
async function fetchGitHubCommits() {
    try {
        const response = await fetch('https://api.andreiixe.website/api/github');
        const commits = await response.json();

        const ul = document.getElementById('github-commits');
        ul.innerHTML = '';

        commits.forEach(commit => {
            const li = document.createElement('li');
            li.innerHTML = `
                <div>
                    <strong>${commit.repo}</strong><br>
                    <a href="${commit.url}" target="_blank">${commit.message}</a>
                </div>
            `;
            ul.appendChild(li);
        });
    } catch (err) {
        console.error('GitHub fetch error:', err);
        document.getElementById('github-commits').innerText = 'Unable to fetch commits.';
    }
}
// Fetch imediat și la fiecare 60 sec
fetchGitHubCommits();
setInterval(fetchGitHubCommits, 1800000); // refresh la 30 minute

// --- Fetch ultimele melodii Last.fm ---
async function fetchLastPlaying() {
    try {
        const response = await fetch('https://api.andreiixe.website/api/lastfm');
        const tracks = await response.json();

        const container = document.getElementById('last-playing');
        container.innerHTML = '';

        let nowPlayingFound = false;

        tracks.forEach(track => {
            const isNowPlaying = track.now_playing === true || track.now_playing === "true";
            let trackHTML = '';

            if (isNowPlaying) {
                trackHTML = `
                    <div class='lead mb-4 text-center'>
                        <img src='src/gifs/dance.gif'><br>
                        <p>⋆.˚✮🎧✮˚.⋆</p>
                        Listening now: <strong class='colors'>${track.name}</strong> 
                        <em class='colors'>by ${track.artist}</em>
                    </div>
                `;
                nowPlayingFound = true;
            } else {
                trackHTML = `
                    <div>⋗ . ♬ ݁˖ 『 <strong>${track.name}</strong> by ${track.artist} 』</div>
                `;
            }

            container.innerHTML += trackHTML;
        });

        if (!nowPlayingFound && tracks.length > 0) {
            container.innerHTML = `<em>Recently played songs:</em>` + container.innerHTML;
        }

    } catch (err) {
        console.error('Last.fm fetch error:', err);
        document.getElementById('last-playing').innerText = 'Unable to fetch Last.fm data.';
    }
}


// Fetch imediat și la fiecare 10 sec
fetchLastPlaying();
setInterval(fetchLastPlaying, 10000);

