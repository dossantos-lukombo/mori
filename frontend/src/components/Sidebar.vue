<template>
    <div :class="['sidebar', { 'sidebar--active': isActive }]" @click.self="closeSidebar">
      <div class="sidebar-content">
        <ul>
          <li><router-link to="/dashboard">Dashboard</router-link></li>
          <li><router-link to="/settings">Settings</router-link></li>
          <li @click="logout">Log out</li>
        </ul>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: "Sidebar",
    props: {
      isActive: {
        type: Boolean,
        required: true,
      },
    },
    methods: {
      closeSidebar() {
        this.$emit('close-sidebar');
      },
      logout() {
        this.$emit('logout');
      }
    },
  };
  </script>
  
  <style scoped>
  .sidebar {
    position: fixed;
    top: 64.45px; /* Adjust to match your navbar height */
    left: -440px;
    width: 440px;
    height: calc(100% - 64.45px); /* Subtract navbar height from full height */
    background-color: var(--bg-neutral);
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3); /* Adjust transparency for softer shadow */
    color: var(--color-white);
    transition: left 0.3s ease;
    border-radius: 0 0 20px 0; /* Rounded bottom corners */
    z-index: 2;
    overflow-y: auto;
    clip-path: inset(-10px -10px 0 -10px); /* Remove shadow from the top */
}

  .sidebar--active {
    left: 0;
  }
  .sidebar-content {
    padding: 20px;
  }
  .sidebar ul {
    list-style: none;
    padding: 0;
  }
  .sidebar li {
    margin: 15px 0;
  }
  .sidebar a {
    color: var(--color-white);
    text-decoration: none;
  }
  .sidebar a:hover {
    text-decoration: underline;
  }
  </style>
  