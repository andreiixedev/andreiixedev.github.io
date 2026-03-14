// Format date function
function formatDate(dateString) {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(dateString).toLocaleDateString('en-US', options);
}

// Get icon class based on type
function getIconClass(type) {
    const icons = {
        'added': 'plus-circle',
        'changed': 'edit',
        'fixed': 'bug',
        'improved': 'arrow-up',
        'removed': 'trash'
    };
    return icons[type] || 'circle';
}

// Load changelog from JSON
async function loadChangelog() {
    try {
        const response = await fetch('changelog.json');
        const data = await response.json();
        
        updateStats(data.stats);
        generateTimeline(data.entries);
        generateFutureFeatures(data.futureFeatures);
        initFilters();
        
    } catch (error) {
        console.error('Error loading changelog.json:', error);
        document.getElementById('changelogTimeline').innerHTML = `
            <div style="text-align: center; padding: 40px; color: #ff4444;">
                <i class="fas fa-exclamation-triangle" style="font-size: 48px; margin-bottom: 15px;"></i>
                <h3>Failed to load changelog data</h3>
                <p style="margin-top: 10px; color: var(--text-muted);">Please check if changelog.json exists</p>
            </div>
        `;
    }
}

// Update statistics
function updateStats(stats) {
    document.getElementById('totalUpdates').textContent = stats.totalUpdates;
    document.getElementById('featuresAdded').textContent = stats.featuresAdded;
    document.getElementById('bugsFixed').textContent = stats.bugsFixed;
    document.getElementById('improvements').textContent = stats.improvements;
}

// Generate timeline entries
function generateTimeline(entries) {
    const timeline = document.getElementById('changelogTimeline');
    timeline.innerHTML = '';
    
    entries.forEach((entry, index) => {
        const changesHTML = entry.changes.map(change => `
            <li>
                <span class="changelog-tag tag-${change.type}">${change.type}</span>
                <i class="fas fa-${change.icon}" style="color: var(--accent); width: 20px;"></i>
                ${change.description}
            </li>
        `).join('');
        
        const itemHTML = `
            <div class="changelog-item" data-version="${entry.version}">
                <div class="changelog-marker">
                    <i class="fas fa-${entry.icon}"></i>
                </div>
                <div class="changelog-content">
                    <div class="changelog-date">
                        <i class="far fa-calendar-alt" style="margin-right: 5px;"></i>
                        ${formatDate(entry.date)}
                        <span class="changelog-version">${entry.version}</span>
                    </div>
                    <h3 class="changelog-title">
                        ${entry.title}
                    </h3>
                    <ul class="changelog-list">
                        ${changesHTML}
                    </ul>
                </div>
            </div>
        `;
        
        timeline.innerHTML += itemHTML;
    });
}

// Generate future features
function generateFutureFeatures(features) {
    const container = document.getElementById('futureFeatures');
    container.innerHTML = features.map(feature => `
        <span class="tag">
            <i class="fas fa-${feature.icon}"></i> ${feature.name}
        </span>
    `).join('');
}

// Initialize filters
function initFilters() {
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            
            const filter = btn.getAttribute('data-filter');
            const items = document.querySelectorAll('.changelog-item');
            
            items.forEach(item => {
                if (filter === 'all') {
                    item.style.display = 'block';
                } else {
                    const hasTag = Array.from(item.querySelectorAll('.changelog-tag')).some(
                        tag => tag.classList.contains(`tag-${filter}`)
                    );
                    item.style.display = hasTag ? 'block' : 'none';
                }
            });
        });
    });
}

// Load changelog when page loads
document.addEventListener('DOMContentLoaded', loadChangelog);