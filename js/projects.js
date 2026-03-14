// Load projects from JSON
async function loadProjects() {
    try {
        const response = await fetch('projects.json');
        const data = await response.json();
        
        // Save data globally
        window.projectsData = data;
        
        // Generate projects HTML
        generateProjects(data.projects);
        
        // Generate category filters
        generateCategories(data.categories);
        
        // Add event listeners
        addProjectListeners();
        
    } catch (error) {
        console.error('Error loading projects:', error);
        document.getElementById('projectsGrid').innerHTML = `
            <div style="grid-column: 1/-1; text-align: center; padding: 40px; color: #ff4444;">
                <i class="fas fa-exclamation-triangle" style="font-size: 48px; margin-bottom: 15px;"></i>
                <h3>Failed to load projects</h3>
                <p style="margin-top: 10px; color: var(--text-muted);">Please check if projects.json exists</p>
            </div>
        `;
    }
}

// Generate projects HTML
function generateProjects(projects) {
    const grid = document.getElementById('projectsGrid');
    grid.innerHTML = '';
    
    projects.forEach(project => {
        // Create technologies tags
        const techTags = project.technologies.map(tech => 
            `<span class="project-tag">${tech}</span>`
        ).join('');
        
        // Generate links HTML based on what's available
        const linksHTML = generateProjectLinks(project.links);
        
        const projectHTML = `
            <div class="project-item" data-id="${project.id}" data-category="${project.category}" data-featured="${project.featured}">
                <img src="${project.image}" alt="${project.title}" class="project-image" loading="lazy">
                <h3>${project.title}</h3>
                <p>${project.description}</p>
                <div class="project-tags">
                    ${techTags}
                </div>
                <div class="project-meta">
                    <span><i class="far fa-calendar-alt"></i> ${project.year}</span>
                    <div class="project-links">
                        ${linksHTML}
                    </div>
                </div>
            </div>
        `;
        
        grid.innerHTML += projectHTML;
    });
}

// Generate links HTML based on available links
function generateProjectLinks(links) {
    let html = '';
    
    // GitHub link (always show if exists)
    if (links.github) {
        html += `<a href="${links.github}" target="_blank" class="project-link" title="GitHub"><i class="fab fa-github"></i></a>`;
    }
    
    // Live demo link (if exists)
    if (links.live) {
        html += `<a href="${links.live}" target="_blank" class="project-link" title="Live Demo"><i class="fas fa-external-link-alt"></i></a>`;
    }
    
    // Release link (if exists and no live demo)
    if (links.release && !links.live) {
        html += `<a href="${links.release}" target="_blank" class="project-link" title="Release"><i class="fas fa-tag"></i></a>`;
    }
    
    // If no links at all, show a disabled state or nothing
    if (!links.github && !links.live && !links.release) {
        html = `<span class="project-link disabled" title="No links available"><i class="fas fa-ban"></i></span>`;
    }
    
    return html;
}

// Generate category filters
function generateCategories(categories) {
    const filterContainer = document.getElementById('categoryFilters');
    if (!filterContainer) return;
    
    filterContainer.innerHTML = '';
    
    categories.forEach(category => {
        const btn = document.createElement('button');
        btn.className = `filter-btn ${category === 'all' ? 'active' : ''}`;
        btn.setAttribute('data-category', category);
        btn.innerHTML = category.charAt(0).toUpperCase() + category.slice(1);
        
        btn.addEventListener('click', () => {
            // Update active button
            document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            
            // Filter projects
            filterProjects(category);
        });
        
        filterContainer.appendChild(btn);
    });
}

// Filter projects by category
function filterProjects(category) {
    const projects = document.querySelectorAll('.project-item');
    const showFeaturedOnly = document.getElementById('featuredToggle')?.checked || false;
    
    projects.forEach(project => {
        const projectCategory = project.dataset.category;
        const isFeatured = project.dataset.featured === 'true';
        
        let show = category === 'all' || projectCategory === category;
        
        if (showFeaturedOnly) {
            show = show && isFeatured;
        }
        
        project.style.display = show ? 'block' : 'none';
    });
    
    // Show "no results" message
    const visibleProjects = document.querySelectorAll('.project-item[style="display: block;"]').length;
    const noResultsMsg = document.getElementById('noResults');
    
    if (noResultsMsg) {
        noResultsMsg.style.display = visibleProjects === 0 ? 'block' : 'none';
    }
}

// Add event listeners for interactive features
function addProjectListeners() {
    // Search functionality
    const searchInput = document.getElementById('searchProjects');
    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            const searchTerm = e.target.value.toLowerCase();
            const projects = document.querySelectorAll('.project-item');
            
            projects.forEach(project => {
                if (project.style.display !== 'none') {
                    const title = project.querySelector('h3').textContent.toLowerCase();
                    const desc = project.querySelector('p').textContent.toLowerCase();
                    const matches = title.includes(searchTerm) || desc.includes(searchTerm);
                    project.style.display = matches ? 'block' : 'none';
                }
            });
            
            // Check for no results
            const visibleProjects = document.querySelectorAll('.project-item[style="display: block;"]').length;
            const noResultsMsg = document.getElementById('noResults');
            if (noResultsMsg) {
                noResultsMsg.style.display = visibleProjects === 0 ? 'block' : 'none';
            }
        });
    }
    
    // Featured toggle
    const featuredToggle = document.getElementById('featuredToggle');
    if (featuredToggle) {
        featuredToggle.addEventListener('change', () => {
            const activeCategory = document.querySelector('.filter-btn.active')?.dataset.category || 'all';
            filterProjects(activeCategory);
        });
    }
    
    // Click on project to show details (but not when clicking on links)
    document.querySelectorAll('.project-item').forEach(project => {
        project.addEventListener('click', (e) => {
            // Don't open if clicking on a link
            if (e.target.tagName === 'A' || e.target.closest('a')) return;
            
            const projectId = project.dataset.id;
            showProjectDetails(projectId);
        });
    });
}

// Show project details in modal
function showProjectDetails(projectId) {
    const project = window.projectsData.projects.find(p => p.id == projectId);
    if (!project) return;
    
    const modal = document.getElementById('projectModal');
    const modalContent = document.getElementById('modalContent');
    
    const techTags = project.technologies.map(tech => 
        `<span class="project-tag">${tech}</span>`
    ).join('');
    
    // Generate links for modal
    const linksHTML = generateModalLinks(project.links);
    
    modalContent.innerHTML = `
        <div style="display: flex; gap: 30px; flex-wrap: wrap;">
            <img src="${project.image}" alt="${project.title}" style="width: 300px; border-radius: 15px; border: 3px solid var(--accent);">
            <div style="flex: 1;">
                <h2 style="font-family: 'DM Serif Display'; color: var(--accent); margin-bottom: 15px;">${project.title}</h2>
                <p style="margin-bottom: 20px;">${project.description}</p>
                <div style="margin-bottom: 20px;">
                    <strong>Technologies:</strong>
                    <div style="margin-top: 10px; display: flex; flex-wrap: wrap; gap: 8px;">${techTags}</div>
                </div>
                <p><strong>Year:</strong> ${project.year}</p>
                <p><strong>Category:</strong> ${project.category}</p>
                <div style="margin-top: 30px; display: flex; gap: 15px; flex-wrap: wrap;">
                    ${linksHTML}
                </div>
            </div>
        </div>
    `;
    
    modal.style.display = 'flex';
    document.body.style.overflow = 'hidden';
}

// Generate links for modal
function generateModalLinks(links) {
    let html = '';
    
    if (links.github) {
        html += `<a href="${links.github}" target="_blank" class="nav-link" style="padding: 10px 25px;"><i class="fab fa-github"></i> GitHub</a>`;
    }
    
    if (links.live) {
        html += `<a href="${links.live}" target="_blank" class="nav-link" style="padding: 10px 25px;"><i class="fas fa-external-link-alt"></i> Live Demo</a>`;
    }
    
    if (links.release && !links.live) {
        html += `<a href="${links.release}" target="_blank" class="nav-link" style="padding: 10px 25px;"><i class="fas fa-tag"></i> Release</a>`;
    }
    
    if (!links.github && !links.live && !links.release) {
        html = `<span class="nav-link disabled" style="opacity: 0.5; cursor: not-allowed;">No links available</span>`;
    }
    
    return html;
}

// Close modal
function closeProjectModal() {
    const modal = document.getElementById('projectModal');
    modal.style.display = 'none';
    document.body.style.overflow = '';
}

// Close modal with Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        closeProjectModal();
    }
});

// Load projects when page loads
document.addEventListener('DOMContentLoaded', loadProjects);