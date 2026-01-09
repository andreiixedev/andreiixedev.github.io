const btn = document.getElementById('toggleMenu');
const mobileNavbar = document.getElementById('mobileNavbar');
const menuLinks = document.querySelectorAll('.menu-links a');
const themeBtn = document.getElementById('themeToggle');
const themeBtnMobile = document.getElementById('themeToggleMobile');
const body = document.body;

function toggleMenu() {
    btn.classList.toggle('active');
    mobileNavbar.querySelector('.menu-links').classList.toggle('active');
}

btn.addEventListener('click', toggleMenu);

menuLinks.forEach(link => {
    link.addEventListener('click', () => {
        if (mobileNavbar.querySelector('.menu-links').classList.contains('active')) {
            toggleMenu();
        }
    });
});

function setTheme(light) {
    if (light) {
        body.setAttribute('data-theme', 'light');
        themeBtn.innerHTML = '<i class="fas fa-sun"></i>';
        themeBtnMobile.innerHTML = '<i class="fas fa-sun"></i>';
        localStorage.setItem("theme", "light");
    } else {
        body.setAttribute('data-theme', 'dark');
        themeBtn.innerHTML = '<i class="fas fa-moon"></i>';
        themeBtnMobile.innerHTML = '<i class="fas fa-moon"></i>';
        localStorage.setItem("theme", "dark");
    }
}

function toggleTheme() {
    const isLight = body.getAttribute('data-theme') === 'light';
    setTheme(!isLight);
}

themeBtn.addEventListener('click', toggleTheme);
themeBtnMobile.addEventListener('click', toggleTheme);

const savedTheme = localStorage.getItem("theme");
if (savedTheme) {
    setTheme(savedTheme === "light");
}

if (document.querySelector(".auto-input")) {
    const texts = ["Andreiixe", "Macintosh", "Pikachu", ":)", "hell"];
    const span = document.querySelector(".auto-input");
    let index = 0;

    function showNextText() {
        span.style.opacity = 0;
        span.style.transform = "translateY(20px)";

        setTimeout(() => {
            span.textContent = texts[index];
            
            span.style.opacity = 1;
            span.style.transform = "translateY(0)";
            span.style.transition = "all 0.4s cubic-bezier(0.4, 0, 0.2, 1)";

            index = (index + 1) % texts.length;
            
            setTimeout(showNextText, 2500);
        }, 400);
    }

    showNextText();
}

// GitHub
async function fetchGitHubCommits() {
    try {
        const response = await fetch('https://api.andreiixe.website/api/github');
        const commit = await response.json();

        const container = document.getElementById('github-commits');
        container.innerHTML = '';

        if (!commit || !commit.repo || !commit.message) {
            container.innerHTML = '<li class="text-muted">No recent commits available.</li>';
            return;
        }

        const commitDate = new Date(commit.date);
        const dateString = commitDate.toLocaleString('en-US', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });


        const li = document.createElement('li');
        li.className = 'commit-item';
        li.innerHTML = `
            <div class="commit-content">
                <div class="commit-header">
                    <span class="commit-repo">${commit.repo}</span>
                    ${commit.is_new ? '<span class="commit-badge">NEW</span>' : ''}
                </div>
                <a href="${commit.url}" target="_blank" class="commit-message">${commit.message}</a>
                <div class="commit-footer">
                    <span class="commit-date">${dateString}</span>
                    <span class="commit-arrow"><i class="fas fa-external-link-alt"></i></span>
                </div>
            </div>
        `;

        container.appendChild(li);

    } catch (err) {
        console.error('GitHub fetch error:', err);
        document.getElementById('github-commits').innerHTML = 
            '<li class="text-muted">Unable to fetch latest commit.</li>';
    }
}

fetchGitHubCommits();
setInterval(fetchGitHubCommits, 1800000);

// Last.fm
async function fetchLastPlaying() {
    try {
        const response = await fetch('https://api.andreiixe.website/api/lastfm'); // TO DO: plays, music time ago stuff
        const tracks = await response.json();

        const container = document.getElementById('last-playing');
        container.innerHTML = '';

        if (!tracks || tracks.length === 0) {
            container.innerHTML = '<div class="no-tracks">No recent tracks found.</div>';
            return;
        }

        let nowPlayingFound = false;

        tracks.forEach(track => {
            const isNowPlaying = track.now_playing === true || track.now_playing === "true";
            
            if (isNowPlaying) {
                // Now playing track
                container.innerHTML = `
                    <div class="now-playing">
                        <div class="album-art">
                            <div class="vinyl">
                                <div class="vinyl-inner">
                                    <div class="vinyl-hole"></div>
                                </div>
                            </div>
                        </div>
                        <div class="track-info">
                            <div class="track-status">NOW PLAYING <span class="live-indicator">● LIVE</span></div>
                            <h4 class="track-name">${track.name}</h4>
                            <p class="track-artist">by ${track.artist}</p>
                            <div class="track-meta">
                                <span class="track-scrobbles">${track.scrobbles || 0} plays</span>
                            </div>
                        </div>
                    </div>
                `;
                nowPlayingFound = true;
            } else {
                // Recent track
                const trackElement = document.createElement('div');
                trackElement.className = 'recent-track';
                trackElement.innerHTML = `
                    <div class="recent-track-content">
                        <i class="fas fa-music"></i>
                        <div class="recent-track-info">
                            <span class="recent-track-name">${track.name}</span>
                            <span class="recent-track-artist">${track.artist}</span>
                        </div>
                        <span class="recent-track-time">${formatTrackTime(track.date)}</span>
                    </div>
                `;
                container.appendChild(trackElement);
            }
        });

        // If no now playing, add a header
        if (!nowPlayingFound && tracks.length > 0) {
            const header = document.createElement('div');
            header.className = 'recent-header';
            header.innerHTML = '<i class="fas fa-history"></i> Recently Played';
            container.insertBefore(header, container.firstChild);
        }

    } catch (err) {
        console.error('Last.fm fetch error:', err);
        document.getElementById('last-playing').innerHTML = 
            '<div class="error-message">Unable to fetch music data.</div>';
    }
}

function formatTrackTime(timestamp) {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 60) {
        return `${diffMins}m ago`;
    } else if (diffHours < 24) {
        return `${diffHours}h ago`;
    } else {
        return `${diffDays}d ago`;
    }
}

fetchLastPlaying();
setInterval(fetchLastPlaying, 10000);

document.addEventListener('DOMContentLoaded', () => {
    const style = document.createElement('style');
    style.textContent = `
        .commit-item {
            padding: 1rem;
            background: var(--glass-bg);
            border-radius: 12px;
            margin-bottom: 0.75rem;
            border: 1px solid var(--glass-border);
            transition: var(--transition);
        }
        
        .commit-item:hover {
            transform: translateX(5px);
            border-color: var(--primary);
        }
        
        .commit-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 0.5rem;
        }
        
        .commit-repo {
            font-weight: 600;
            color: var(--primary);
            font-size: 0.9rem;
        }
        
        .commit-badge {
            background: linear-gradient(135deg, var(--accent), #10b981);
            color: white;
            font-size: 0.7rem;
            padding: 0.2rem 0.5rem;
            border-radius: 4px;
            font-weight: 600;
        }
        
        .commit-message {
            color: var(--dark);
            text-decoration: none;
            font-weight: 500;
            display: block;
            margin-bottom: 0.5rem;
            line-height: 1.4;
        }
        
        .commit-message:hover {
            color: var(--primary);
        }
        
        .commit-footer {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 0.8rem;
            color: var(--gray);
        }
        
        .commit-arrow {
            opacity: 0.6;
            transition: var(--transition);
        }
        
        .commit-item:hover .commit-arrow {
            opacity: 1;
            color: var(--primary);
        }
        
        .music-content {
            padding: 1rem 0;
        }
        
        .now-playing {
            display: flex;
            align-items: center;
            gap: 1.5rem;
            padding: 1.5rem;
            background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(139, 92, 246, 0.1));
            border-radius: 16px;
            border: 1px solid var(--glass-border);
            margin-bottom: 1.5rem;
        }
        
        .album-art {
            position: relative;
            width: 100px;
            height: 100px;
        }
        
        .vinyl {
            width: 100%;
            height: 100%;
            background: linear-gradient(45deg, #1e293b, #0f172a);
            border-radius: 50%;
            position: relative;
            animation: rotate 20s linear infinite;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
        }
        
        @keyframes rotate {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .vinyl-inner {
            position: absolute;
            top: 10%;
            left: 10%;
            right: 10%;
            bottom: 10%;
            background: linear-gradient(45deg, #334155, #475569);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        
        .vinyl-hole {
            width: 10px;
            height: 10px;
            background: var(--dark);
            border-radius: 50%;
        }
        
        .track-info {
            flex: 1;
        }
        
        .track-status {
            font-size: 0.8rem;
            color: var(--accent);
            font-weight: 600;
            margin-bottom: 0.5rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        
        .live-indicator {
            color: #ef4444;
            animation: pulse 1.5s infinite;
        }
        
        .track-name {
            font-size: 1.5rem;
            font-weight: 700;
            margin-bottom: 0.25rem;
            color: var(--dark);
        }
        
        .track-artist {
            color: var(--gray);
            font-size: 1.1rem;
            margin-bottom: 1rem;
        }
        
        .track-meta {
            display: flex;
            gap: 1rem;
            font-size: 0.9rem;
            color: var(--gray);
        }
        
        .track-scrobbles {
            background: var(--glass-bg);
            padding: 0.25rem 0.75rem;
            border-radius: 20px;
            border: 1px solid var(--glass-border);
        }
        
        .recent-header {
            font-size: 0.9rem;
            color: var(--gray);
            margin-bottom: 1rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        
        .recent-track {
            padding: 0.75rem;
            border-radius: 12px;
            margin-bottom: 0.5rem;
            transition: var(--transition);
        }
        
        .recent-track:hover {
            background: var(--glass-bg);
        }
        
        .recent-track-content {
            display: flex;
            align-items: center;
            gap: 1rem;
        }
        
        .recent-track-content i {
            color: var(--primary);
            width: 20px;
        }
        
        .recent-track-info {
            flex: 1;
            display: flex;
            flex-direction: column;
        }
        
        .recent-track-name {
            font-weight: 500;
            color: var(--dark);
            font-size: 0.9rem;
        }
        
        .recent-track-artist {
            font-size: 0.8rem;
            color: var(--gray);
        }
        
        .recent-track-time {
            font-size: 0.8rem;
            color: var(--gray);
        }
        
        .no-tracks, .error-message {
            text-align: center;
            padding: 2rem;
            color: var(--gray);
            font-style: italic;
        }
    `;
    document.head.appendChild(style);
});