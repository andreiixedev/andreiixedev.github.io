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

            span.textContent = texts[index];

 
            span.style.opacity = 1;
            span.style.transform = "translateY(0) scale(1)";


            index = (index + 1) % texts.length;
            setTimeout(showNextText, 2500); 
        }, 400); 
    }

    showNextText();
}



async function fetchGitHubCommits() {
    try {
        const response = await fetch('https://api.andreiixe.website/api/github');
        const commit = await response.json();

        const container = document.getElementById('github-commits');
        container.innerHTML = '';


        if (!commit || !commit.repo || !commit.message) {
            container.innerText = 'Nu există commituri recente.';
            return;
        }


        const commitDate = new Date(commit.date);
        const dateString = commitDate.toLocaleString('ro-RO', {
            day: '2-digit',
            month: 'short',
            hour: '2-digit',
            minute: '2-digit'
        });


        const li = document.createElement('li');
        li.innerHTML = `
            <div style="display: flex; flex-direction: column; gap: 4px;">
                <div>
                    <strong>${commit.repo}</strong>
                    ${commit.is_new ? '<span style="color: limegreen; font-weight: bold;">🟢 New</span>' : ''}
                </div>
                <a href="${commit.url}" target="_blank" style="color: #00bfff;">${commit.message}</a>
                <small style="opacity: 0.7;">${dateString}</small>
            </div>
        `;

        container.appendChild(li);
    } catch (err) {
        console.error('GitHub fetch error:', err);
        document.getElementById('github-commits').innerText = 'Unable to fetch latest commit.';
    }
}

fetchGitHubCommits();
setInterval(fetchGitHubCommits, 1800000);


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


fetchLastPlaying();
setInterval(fetchLastPlaying, 10000);

