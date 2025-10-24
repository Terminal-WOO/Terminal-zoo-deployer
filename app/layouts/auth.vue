<template>
  <div class="auth-shell">
    <div class="auth-frame">
      <div class="auth-grid">
        <!-- LEFT: pagina-content (login/register/forgot) -->
        <section class="auth-left">
          <div class="brand">
            <BrandLogo size="lg" />
            
          </div>

          <!-- de pagina vult dit slot -->
          <slot />
        </section>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/* niets nodig */
</script>

<style scoped>
/* Shell */
.auth-shell {
  min-height: 100vh;
  background: var(--surface-50);
  display: grid;
  place-items: center;
  padding: 2rem;
}

/* Outer frame */
.auth-frame {
  width: 100%;
  max-width: 1200px;
  border-radius: 24px;
  background: var(--surface-0);
  border: 1px solid var(--surface-200);
  box-shadow: 0 30px 80px rgba(0,0,0,.08);
  overflow: hidden;
}

/* Grid */
.auth-grid {
  display: grid;
  grid-template-columns: 1fr 1.15fr; /* ~ 46% / 54% */
}

@media (max-width: 992px) {
  .auth-grid { grid-template-columns: 1fr; }
}

/* Left side */
.auth-left {
  padding: 2.25rem 2.5rem 2.5rem;
}

.brand { display: flex; align-items: center; gap: .75rem; margin-bottom: 1.25rem; }
.brand-name { font-weight: 700; color: var(--ai-deep-green, #0E2D19); }
.brand-sub  { color: var(--text-color-secondary); }

/* Right side */
.auth-right { background: #f7f9fb; position: relative; }

.hero-bg {
  padding: 1.25rem 1.5rem;
  background: linear-gradient(180deg, #eef6f1 0%, #f7f9fb 100%);
  border-bottom: 1px solid var(--surface-200);
}
.hero-title { font-size: .75rem; color: var(--text-color-secondary); }
.hero-sub { margin-top: .25rem; color: var(--text-color); }

/* Dashboard mock */
.mock {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: 16px;
  padding: 16px;
}
.mock-sidebar {
  background: white;
  border: 1px solid var(--surface-200);
  border-radius: 12px;
  padding: 12px;
  display: grid; grid-template-rows: 40px 1fr 40px; gap: 12px;
}
.mock-logo { background: var(--surface-200); height: 40px; border-radius: 8px; }
.mock-nav { list-style: none; padding: 0; margin: 0; display: grid; gap: 10px; }
.mock-nav li { height: 36px; border-radius: 8px; background: var(--surface-200); }
.mock-nav li.active { background: var(--primary-100, #E6F8EB); border: 1px solid var(--primary-200, #99DBAD); }
.mock-footer { background: var(--surface-200); border-radius: 8px; }

.mock-main { display: grid; gap: 16px; }
.mock-cards { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.mock-card { height: 66px; border-radius: 10px; background: white; border: 1px solid var(--surface-200); }

.mock-board {
  background: white; border: 1px solid var(--surface-200);
  border-radius: 12px; padding: 12px; display: grid; gap: 10px;
}
.mock-row { height: 52px; border-radius: 8px; background: var(--surface-100); }
.mock-row.light { background: var(--surface-50); }
</style>
