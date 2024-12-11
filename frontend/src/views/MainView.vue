<template>
    <NavBarOn />
    <div id="layout">
      <ChatbotConversation />
    </div>
  </template>
  
  <script>
  import NavBarOn from '@/components/NavBarOn.vue';
  import ChatbotConversation from '@/components/ChatbotConversation.vue';
  
  export default {
    name: 'MainView',
    components: { NavBarOn, ChatbotConversation },
  };
  </script>
  
  <style>
  /* Ensure the parent takes up the full height of the page */
  html, body, #app {
    height: 100%;
    margin: 0;
    padding: 0;
  }
  
  /* Full-page layout excluding the navbar */
  #layout {
    width: 100%;
    height: calc(100vh - 64px); /* Adjust 64px to match your NavBar height */
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
  </style>
  