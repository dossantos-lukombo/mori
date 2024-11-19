document.addEventListener('DOMContentLoaded', function() {
    const toggleButton = document.getElementById('toggle-button');
    const sidebar = document.getElementById('sidebar');
    const chatContainer = document.querySelector('.chat-container');
    const chatWindow = document.querySelector('.chat-window');
    const user_input = document.querySelector('.chat-input');
  
    toggleButton.addEventListener('click', function() {
        sidebar.classList.toggle('hidden');
        chatContainer.classList.toggle('full-width');
        chatWindow.scrollTop = chatWindow.scrollHeight;
        chatWindow.classList.toggle('full-width');
        user_input.classList.toggle('full-width');

    });
  });
  