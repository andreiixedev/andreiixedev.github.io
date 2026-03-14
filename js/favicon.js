(function() {
    function addAnimatedFireFavicon() {
        if (window.fireFaviconAdded) return;
        window.fireFaviconAdded = true;
        
        const canvas = document.createElement('canvas');
        canvas.width = 64;
        canvas.height = 64;
        const ctx = canvas.getContext('2d');
        
        let particles = [];
        let smokeParticles = [];
        let sparks = [];
        let embers = [];
        let frame = 0;
        
        function drawPixel(x, y, size, color) {
            ctx.fillStyle = color;
            ctx.fillRect(Math.floor(x), Math.floor(y), size, size);
        }
        
        function initParticles() {
            for (let i = 0; i < 35; i++) {
                particles.push({
                    x: 32 + (Math.random() - 0.5) * 20,
                    y: 50 + Math.random() * 10,
                    size: 2 + Math.random() * 3,
                    speedY: -1.5 - Math.random() * 2,
                    speedX: (Math.random() - 0.5) * 1.5,
                    life: 0.7 + Math.random() * 0.3,
                    flicker: Math.random() * 10
                });
            }
            
            for (let i = 0; i < 25; i++) {
                smokeParticles.push({
                    x: 32 + (Math.random() - 0.5) * 15,
                    y: 35 + Math.random() * 15,
                    size: 3 + Math.random() * 5,
                    speedY: -0.8 - Math.random() * 1,
                    speedX: (Math.random() - 0.5) * 0.8,
                    opacity: 0.2 + Math.random() * 0.4,
                    life: 1
                });
            }
            
            for (let i = 0; i < 20; i++) {
                sparks.push({
                    x: 32 + (Math.random() - 0.5) * 25,
                    y: 40 + Math.random() * 15,
                    size: 1 + Math.random() * 2,
                    speedY: -2 - Math.random() * 3,
                    speedX: (Math.random() - 0.5) * 2.5,
                    life: 0.8 + Math.random() * 0.2,
                    color: Math.random() > 0.5 ? '#ffff00' : '#ffaa00'
                });
            }
            
            for (let i = 0; i < 15; i++) {
                embers.push({
                    x: 32 + (Math.random() - 0.5) * 18,
                    y: 48 + Math.random() * 8,
                    size: 1.5 + Math.random() * 2,
                    speedY: -0.5 - Math.random() * 1.2,
                    speedX: (Math.random() - 0.5) * 1,
                    life: 0.9 + Math.random() * 0.1,
                    glow: true
                });
            }
        }
        
        initParticles();
        
        function drawAnimatedFire() {
            ctx.clearRect(0, 0, 64, 64);
            
            ctx.globalAlpha = 0.9;
            ctx.fillStyle = '#4a4a4a';
            for (let x = 24; x < 40; x += 2) {
                for (let y = 52; y < 60; y += 2) {
                    if (Math.random() > 0.4) {
                        ctx.globalAlpha = 0.7 + Math.random() * 0.3;
                        drawPixel(x + (Math.random() - 0.5) * 2, y, 2, '#333333');
                    }
                }
            }
            
            ctx.globalAlpha = 0.95;
            for (let i = 0; i < 8; i++) {
                let x = 28 + Math.sin(frame * 0.1 + i) * 4 + Math.random() * 2;
                let y = 52 + Math.cos(i) * 2;
                ctx.globalAlpha = 0.9;
                drawPixel(x, y, 3, '#ff5500');
                ctx.globalAlpha = 0.8;
                drawPixel(x-1, y-1, 2, '#ffaa00');
            }
            
            smokeParticles.forEach((p, index) => {
                let opacity = p.opacity * p.life * 0.7;
                
                const gradient = ctx.createRadialGradient(p.x, p.y, 0, p.x, p.y, p.size * 1.5);
                gradient.addColorStop(0, `rgba(120, 120, 120, ${opacity})`);
                gradient.addColorStop(0.5, `rgba(90, 90, 90, ${opacity * 0.6})`);
                gradient.addColorStop(1, `rgba(60, 60, 60, 0)`);
                
                ctx.fillStyle = gradient;
                ctx.globalAlpha = 1;
                ctx.beginPath();
                ctx.arc(p.x, p.y, p.size * 1.2, 0, Math.PI * 2);
                ctx.fill();
                
                p.x += p.speedX + Math.sin(frame * 0.1 + index) * 0.3;
                p.y += p.speedY;
                p.size += 0.05;
                p.life -= 0.002;
                p.opacity -= 0.001;
                
                if (p.life <= 0.2 || p.y < 10 || p.x < 5 || p.x > 59) {
                    smokeParticles[index] = {
                        x: 32 + (Math.random() - 0.5) * 15,
                        y: 45 + Math.random() * 10,
                        size: 3 + Math.random() * 4,
                        speedY: -0.6 - Math.random() * 1,
                        speedX: (Math.random() - 0.5) * 0.6,
                        opacity: 0.2 + Math.random() * 0.3,
                        life: 1
                    };
                }
            });
            

            particles.forEach((p, index) => {
                const gradient = ctx.createRadialGradient(p.x, p.y, 0, p.x, p.y, p.size * 2.5);
                
                if (Math.random() > 0.6) {
                    gradient.addColorStop(0, 'rgba(255, 255, 0, 1)');
                    gradient.addColorStop(0.3, 'rgba(255, 170, 0, 0.9)');
                } else {
                    gradient.addColorStop(0, 'rgba(255, 170, 0, 1)');
                    gradient.addColorStop(0.3, 'rgba(255, 85, 0, 0.9)');
                }
                gradient.addColorStop(0.7, 'rgba(255, 50, 0, 0.5)');
                gradient.addColorStop(1, 'rgba(255, 0, 0, 0)');
                
                ctx.fillStyle = gradient;
                ctx.globalAlpha = 1;
                ctx.beginPath();
                ctx.arc(p.x, p.y, p.size * 1.8, 0, Math.PI * 2);
                ctx.fill();
                
                ctx.fillStyle = 'rgba(255, 255, 170, 0.9)';
                ctx.beginPath();
                ctx.arc(p.x - 0.5, p.y - 0.5, p.size * 0.5, 0, Math.PI * 2);
                ctx.fill();
                
                p.x += p.speedX + Math.sin(frame * 0.2 + index) * 0.4;
                p.y += p.speedY;
                p.size += (Math.random() - 0.5) * 0.3;
                if (p.size < 1.5) p.size = 2;
                if (p.size > 4) p.size = 3.5;
                
                if (frame % Math.floor(p.flicker) === 0) {
                    p.speedX += (Math.random() - 0.5) * 0.2;
                }
                
                if (p.y < 20 || p.x < 10 || p.x > 54) {
                    particles[index] = {
                        x: 32 + (Math.random() - 0.5) * 18,
                        y: 52 + Math.random() * 8,
                        size: 2 + Math.random() * 2.5,
                        speedY: -1.5 - Math.random() * 2,
                        speedX: (Math.random() - 0.5) * 1.2,
                        life: 0.7 + Math.random() * 0.3,
                        flicker: Math.random() * 10
                    };
                }
            });
            
            sparks.forEach((s, index) => {
                ctx.shadowColor = 'rgba(255, 170, 0, 0.8)';
                ctx.shadowBlur = 8;
                
                ctx.fillStyle = s.color;
                ctx.globalAlpha = s.life * 0.9;
                ctx.beginPath();
                ctx.arc(s.x, s.y, s.size, 0, Math.PI * 2);
                ctx.fill();
                
                ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
                ctx.globalAlpha = s.life * 0.5;
                ctx.beginPath();
                ctx.arc(s.x - 0.5, s.y - 0.5, s.size * 0.4, 0, Math.PI * 2);
                ctx.fill();
                
                s.x += s.speedX;
                s.y += s.speedY;
                s.size *= 0.99;
                s.life -= 0.005;
                
                if (s.life <= 0.2 || s.y < 10) {
                    sparks[index] = {
                        x: 32 + (Math.random() - 0.5) * 22,
                        y: 48 + Math.random() * 10,
                        size: 1 + Math.random() * 2,
                        speedY: -2.5 - Math.random() * 3,
                        speedX: (Math.random() - 0.5) * 3,
                        life: 0.8 + Math.random() * 0.2,
                        color: Math.random() > 0.5 ? '#ffffaa' : '#ffaa00'
                    };
                }
            });
            
            embers.forEach((e, index) => {
                ctx.shadowBlur = 12;
                ctx.shadowColor = 'rgba(255, 85, 0, 0.6)';
                
                ctx.fillStyle = 'rgba(255, 68, 0, 0.9)';
                ctx.globalAlpha = e.life * 0.8;
                ctx.beginPath();
                ctx.arc(e.x, e.y, e.size, 0, Math.PI * 2);
                ctx.fill();
                
                ctx.fillStyle = 'rgba(255, 136, 0, 0.8)';
                ctx.beginPath();
                ctx.arc(e.x - 0.5, e.y - 0.5, e.size * 0.6, 0, Math.PI * 2);
                ctx.fill();
                
                e.x += e.speedX + Math.sin(frame * 0.3 + index) * 0.2;
                e.y += e.speedY;
                e.life -= 0.001;
                
                if (e.life <= 0.3 || e.y < 25) {
                    embers[index] = {
                        x: 32 + (Math.random() - 0.5) * 16,
                        y: 50 + Math.random() * 6,
                        size: 1.5 + Math.random() * 2,
                        speedY: -0.8 - Math.random() * 1.5,
                        speedX: (Math.random() - 0.5) * 1.2,
                        life: 0.9 + Math.random() * 0.1,
                        glow: true
                    };
                }
            });
            
            ctx.shadowBlur = 15;
            ctx.shadowColor = 'rgba(255, 100, 0, 0.3)';
            ctx.fillStyle = 'rgba(255, 100, 0, 0.03)';
            ctx.beginPath();
            ctx.arc(32, 45, 18, 0, Math.PI * 2);
            ctx.fill();
            
            ctx.shadowBlur = 0;
            ctx.globalAlpha = 1;
            
            frame++;
            
            const smallCanvas = document.createElement('canvas');
            smallCanvas.width = 32;
            smallCanvas.height = 32;
            const smallCtx = smallCanvas.getContext('2d');
            
            smallCtx.clearRect(0, 0, 32, 32);
            smallCtx.imageSmoothingEnabled = false;
            smallCtx.drawImage(canvas, 0, 0, 64, 64, 0, 0, 32, 32);
            
            const link = document.querySelector("link[rel*='icon']") || document.createElement('link');
            link.type = 'image/x-icon';
            link.rel = 'shortcut icon';
            link.href = smallCanvas.toDataURL('image/x-icon');
            document.getElementsByTagName('head')[0].appendChild(link);
            
            requestAnimationFrame(drawAnimatedFire);
        }
        
        drawAnimatedFire();
    }
    
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', addAnimatedFireFavicon);
    } else {
        addAnimatedFireFavicon();
    }
})();