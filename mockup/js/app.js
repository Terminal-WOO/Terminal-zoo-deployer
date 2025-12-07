// Module data
const modules = [
    {
        id: '1.1',
        title: 'Platform Engineering Fundamentals',
        phase: 'Foundation',
        description: 'Visie, principes, roadmap en architecture decisions',
        status: 'completed',
        features: [
            'Platform visie en missie',
            'Platform product domains',
            'Strategische roadmap',
            '5 Architecture Decision Records'
        ]
    },
    {
        id: '1.2',
        title: 'Software-Defined Platform Architecture',
        phase: 'Foundation',
        description: 'SDLC, fitness functions en CI/CD pipeline',
        status: 'completed',
        features: [
            'Platform SDLC (5 fasen)',
            'Architectural fitness functions',
            'Platform CI/CD pipeline',
            'Domain boundaries checker'
        ]
    },
    {
        id: '1.3',
        title: 'Metrics & Measurement Framework',
        phase: 'Foundation',
        description: 'DORA metrics, platform value metrics, cognitive load',
        status: 'completed',
        features: [
            'DORA metrics collectors',
            'Prometheus metrics exporter',
            'Grafana dashboards',
            'Developer sentiment survey'
        ]
    },
    {
        id: '2.1',
        title: 'Governance, Compliance & Trust',
        phase: 'Building',
        description: 'Policy-as-code, zero-trust, supply chain security',
        status: 'completed',
        features: [
            'OPA policies',
            'Zero-trust network policies',
            'Supply chain security',
            'Developer autonomy framework'
        ]
    },
    {
        id: '2.2',
        title: 'Evolutionary Observability Platform',
        phase: 'Building',
        description: 'Observability als service, SLOs als code',
        status: 'completed',
        features: [
            'SLOs als code (3 SLOs)',
            'Single pane of glass dashboard',
            'Observability hooks',
            'Prometheus/Loki/Tempo stack'
        ]
    },
    {
        id: '2.3',
        title: 'Software-Defined Infrastructure Platform',
        phase: 'Building',
        description: 'Infrastructure as code, pipelines, GitOps',
        status: 'completed',
        features: [
            'Infrastructure pipelines',
            'Manifest validation tests',
            'Infrastructure CI/CD',
            'Terraform modules'
        ]
    },
    {
        id: '2.4',
        title: 'Platform Control Plane Foundations',
        phase: 'Building',
        description: 'Account baseline, network, identity, Kubernetes',
        status: 'in-progress',
        features: [
            'Cloud account baseline',
            'Transit network layer',
            'Customer identity (OIDC)',
            'Kubernetes control plane'
        ]
    },
    {
        id: '2.5',
        title: 'Control Plane Services & Extensions',
        phase: 'Building',
        description: 'Storage, autoscaling, service mesh, APIs',
        status: 'in-progress',
        features: [
            'Kubernetes storage classes',
            'Cluster autoscaling',
            'Service mesh (Istio/Linkerd)',
            'Platform management APIs'
        ]
    },
    {
        id: '3.1',
        title: 'Architecture for Scale',
        phase: 'Scaling',
        description: 'Event-driven automation, federated control planes',
        status: 'completed',
        features: [
            'Event-driven automation',
            'Federated control planes',
            'Adapter pattern',
            'Distributed orchestration'
        ]
    },
    {
        id: '3.2',
        title: 'Platform Product Evolution',
        phase: 'Scaling',
        description: 'Product mindset, DevEx, intelligent tools, IDP',
        status: 'completed',
        features: [
            'Product evolution strategie',
            'Developer experience',
            'Culture principles',
            'Intelligent tools structuur'
        ]
    }
];

// Render modules
function renderModules() {
    const grid = document.getElementById('modulesGrid');
    if (!grid) return;

    grid.innerHTML = modules.map(module => `
        <div class="module-card ${module.status}">
            <div class="module-card-header">
                <div>
                    <div class="module-card-phase">${module.phase}</div>
                    <div class="module-card-title">Module ${module.id}: ${module.title}</div>
                </div>
            </div>
            <div class="module-card-description">${module.description}</div>
            <ul class="module-card-features">
                ${module.features.map(feature => `
                    <li><i class="fas fa-check"></i> ${feature}</li>
                `).join('')}
            </ul>
        </div>
    `).join('');
}

// Smooth scroll
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Open documentation files
function openDoc(path) {
    // GitHub repo URL
    const githubRepo = 'Terminal-WOO/Terminal-zoo-deployer';
    const githubUrl = `https://github.com/${githubRepo}/blob/main/${path}`;
    
    // Try to open via GitHub (works everywhere)
    window.open(githubUrl, '_blank');
    
    // Also try to open locally if running from file system
    // This will only work if opened from the project directory
    try {
        const localPath = `../${path}`;
        // Only try local if we're on file:// protocol
        if (window.location.protocol === 'file:') {
            // Local file won't work in browser, so just use GitHub
            console.log('Opening via GitHub:', githubUrl);
        }
    } catch (e) {
        // Ignore errors, GitHub link should work
    }
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    renderModules();
    
    // Add active state to nav links
    const navLinks = document.querySelectorAll('.nav-link');
    const sections = document.querySelectorAll('.section');
    
    function updateActiveNav() {
        let current = '';
        sections.forEach(section => {
            const sectionTop = section.offsetTop;
            const sectionHeight = section.clientHeight;
            if (window.scrollY >= sectionTop - 200) {
                current = section.getAttribute('id');
            }
        });
        
        navLinks.forEach(link => {
            link.classList.remove('active');
            if (link.getAttribute('href') === `#${current}`) {
                link.classList.add('active');
            }
        });
    }
    
    window.addEventListener('scroll', updateActiveNav);
    updateActiveNav();
});

