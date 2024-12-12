<template>
    <div id="navbar">
      <div id="menu-btn" @click="toggleSidebar"></div>
      <div id="nav-titleSearch">
        <div class="smallMoriLogo">
        </div>
        <router-link to="/main" class="mori" id="nav-title">Mori</router-link>
        <Search />
      </div>
      <ul class="nav-links">
        <li id="notifications-link">
          <Notifications />
        </li>
        <li>
          <router-link v-if="typeof user.id !== 'undefined'"
                       :to="{ name: 'Profile', params: { id: user.id } }">My profile</router-link>
        </li>
        <li @click="logout">Log out</li>
      </ul>
      <Sidebar :isActive="isSidebarActive" @close-sidebar="toggleSidebar" @logout="logout" />
    </div>
</template>
  
<script>
  import Search from './Search.vue';
  import Notifications from './Notifications.vue';
  import Sidebar from './Sidebar.vue';
  
  export default {
    name: 'NavBarOn',
    data() {
      return {
        user: {},
        isSidebarActive: false,
      };
    },
    components: { Notifications, Search, Sidebar },
    created() {
      this.getUserInfo();
    },
    methods: {
      async getUserInfo() {
        await fetch("http://localhost:8081/currentUser", { credentials: 'include' })
          .then(r => r.json())
          .then(json => {
            this.user = json.users[0];
          });
      },
      async logout() {
        await fetch('http://localhost:8081/logout', {
          credentials: 'include',
          headers: { 'Accept': 'application/json' },
        })
          .then(response => response.json())
          .then(json => console.log(json));
        this.$store.state.wsConn.close(1000, "user logged out");
        this.$router.push("/");
      },
      toggleSidebar() {
        this.isSidebarActive = !this.isSidebarActive;
      }
    }
  };
</script>
  

<style scoped>

#menu-btn {
    width: 30px;
    height: 30px;
    margin-right: 20px;
    background-image: url('../assets/menu.png');
    background-size: cover;
    cursor: pointer;
    transition: all 0.3s;
}

#menu-btn:hover {
    transform: scale(1.05);
}

#navbar {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 3;
    width: 100%;
    min-width: min-content;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 30px;
    background-color: var(--purple-color);
    color: var(--color-white);
    position: sticky;
}


#navbar a {
    color: var(--color-white);
}


#nav-title {
    user-select: none;
    position: relative;
}


.nav-links li {
    user-select: none;
    font-weight: 300;
    display: inline-block;
    margin-left: 20px;
    cursor: pointer;

    position: relative;
}


#nav-titleSearch {
    display: flex;
    gap: 25px;
    flex-grow: 1;
    align-items: center;


}


#navbar li:not(#notifications-link)::after,
#nav-title::after {
    content: "";
    height: 2.5px;
    width: 0;
    display: block;
    position: absolute;

    transition: all 0.35s ease-out;
}

#navbar li:not(#notifications-link):hover::after,
#nav-title:hover::after {
    width: 100%;
    background-color: var(--hover-background-color);
}

a:link {
    text-decoration: none;
}

a:visited {
    text-decoration: none;
}

a:hover {
    text-decoration: none;
}

a:active {
    text-decoration: none;
}
</style>