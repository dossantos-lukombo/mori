<template>
  <div class="container">
    <!-- Title -->
    <div class="small_logo"></div>
    <div class="title">Mori</div>

    <!-- Login Form -->
    <form v-if="activeForm === 'login'" @submit.prevent="submitLogin">
      <input
        type="text"
        name="username_email"
        placeholder="Username or Email"
        v-model="loginData.username_email"
        required
      />
      <input
        type="password"
        name="password"
        placeholder="Password"
        v-model="loginData.password"
        required
      />
      <button type="submit">Connexion</button>
      <button type="button" @click="toggleForms('register')">Register</button>
      <div v-if="loginError" class="error-message">{{ loginError }}</div>
      <p id="forgot_password" @click="toggleForms('reset')">
        Forgot your password ?
      </p>
    </form>

    <!-- Register Form -->
    <form v-if="activeForm === 'register'" @submit.prevent="submitRegister">
      <input
        type="text"
        name="username"
        placeholder="Username"
        v-model="registerData.username"
        required
      />
      <input
        type="email"
        name="email"
        placeholder="Email"
        v-model="registerData.email"
        required
      />
      <input
        type="password"
        name="password"
        placeholder="Password"
        v-model="registerData.password"
        required
      />
      <input
        type="password"
        name="confirm_password"
        placeholder="Confirm Password"
        v-model="registerData.confirm_password"
        required
      />
      <img
        id="captcha-image"
        :src="captchaSrc"
        alt="Captcha"
        @click="reloadCaptcha"
      />
      <input
        type="text"
        name="captcha_input"
        placeholder="Enter Captcha"
        v-model="registerData.captcha_input"
        required
      />

      <button type="submit">Register</button>
      <button type="button" @click="toggleForms('login')">Login</button>
      <div v-if="registerError" class="error-message">{{ registerError }}</div>
    </form>

    <!-- Reset Password Form -->
    <form v-if="activeForm === 'reset'" @submit.prevent="submitResetPassword">
      <p class="info_form">A link will be sent to your Email inbox.</p>
      <input
        type="email"
        name="email"
        placeholder="Email"
        v-model="resetPasswordData.email"
        required
      />
      <button type="submit">Send Email</button>
      <button type="button" @click="toggleForms('login')">Login</button>
      <div v-if="resetPasswordError" class="error-message">
        {{ resetPasswordError }}
      </div>
    </form>

    <!-- Popup Notification -->
    <div id="popup" :class="['popup', { hidden: !showPopup }]">
      <span id="popup-message">{{ popupMessage }}</span>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      activeForm: 'login',
      loginData: {
        username_email: '',
        password: '',
      },
      registerData: {
        username: '',
        email: '',
        password: '',
        confirm_password: '',
        captcha_input: '',
      },
      resetPasswordData: {
        email: '',
      },
      loginError: '',
      registerError: '',
      resetPasswordError: '',
      showPopup: false,
      popupMessage: '',
      captchaSrc: '/captcha',
    };
  },
  methods: {
    // Fonction pour récupérer un cookie par nom
    getCookie(name) {
      const cookies = document.cookie.split('; ');
      for (let i = 0; i < cookies.length; i++) {
        const [key, value] = cookies[i].split('=');
        if (key === name) return decodeURIComponent(value);
      }
      return null;
    },
    toggleForms(form) {
      this.activeForm = form;
      this.clearErrors();
    },
    clearErrors() {
      this.loginError = '';
      this.registerError = '';
      this.resetPasswordError = '';
    },
    showPopupMessage(message) {
      this.popupMessage = message;
      this.showPopup = true;
      setTimeout(() => {
        this.showPopup = false;
      }, 3000);
    },
    reloadCaptcha() {
      this.captchaSrc = `/captcha?${new Date().getTime()}`;
    },
    async submitLogin() {
      this.loginError = '';

      const csrfToken = this.getCookie('csrf_token');

      if (!csrfToken) {
        this.loginError = 'CSRF token is missing.';
        return;
      }

      const formData = new FormData();
      formData.append('username_email', this.loginData.username_email);
      formData.append('password', this.loginData.password);
      formData.append('csrf_token', csrfToken);

      try {
        const response = await fetch('/login', {
          method: 'POST',
          body: formData,
        });

        if (!response.ok) {
          const result = await response.json();
          this.loginError =
            result.error || 'An error occurred during login.';
        } else {
          window.location.href = response.url || '/dashboard';
        }
      } catch (error) {
        this.loginError = 'An error occurred during login.';
        console.error(error);
      }
    },
    async submitRegister() {
      this.registerError = '';

      const csrfToken = this.getCookie('csrf_token');

      if (!csrfToken) {
        this.registerError = 'CSRF token is missing.';
        return;
      }

      const formData = new FormData();
      formData.append('username', this.registerData.username);
      formData.append('email', this.registerData.email);
      formData.append('password', this.registerData.password);
      formData.append('confirm_password', this.registerData.confirm_password);
      formData.append('captcha_input', this.registerData.captcha_input);
      formData.append('csrf_token', csrfToken);

      try {
        const response = await fetch('/register', {
          method: 'POST',
          body: formData,
        });

        const result = await response.json();
        if (!response.ok) {
          this.registerError =
            result.error || 'An error occurred during registration.';
          if (result.reloadCaptcha) {
            this.reloadCaptcha();
          }
        } else {
          this.showPopupMessage(
            'Registration successful, please verify your email.'
          );
          this.toggleForms('login');
        }
      } catch (error) {
        this.registerError = 'An error occurred during registration.';
        console.error(error);
      }
    },
    async submitResetPassword() {
      this.resetPasswordError = '';

      const csrfToken = this.getCookie('csrf_token');

      if (!csrfToken) {
        this.resetPasswordError = 'CSRF token is missing.';
        return;
      }

      const formData = new FormData();
      formData.append('email', this.resetPasswordData.email);
      formData.append('csrf_token', csrfToken);

      try {
        const response = await fetch('/reset-password', {
          method: 'POST',
          body: formData,
        });

        const result = await response.json();
        if (!response.ok) {
          this.resetPasswordError =
            result.error || 'An error occurred during password reset.';
        } else {
          this.showPopupMessage('Password reset link sent to your Email.');
          this.toggleForms('login');
        }
      } catch (error) {
        this.resetPasswordError = 'An error occurred during password reset.';
        console.error(error);
      }
    },
  },
};
</script>


<style>
@font-face {
  font-family: 'TITLE';
  src: url(../assets/fonts/SinosansRegular-aYxZ5.otf) format('truetype');
  font-weight: normal;
  font-style: normal;
}

/* General Body Styling */
body {
  font-family: 'Arial', sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  margin: 0;
  background-color: #1a1a1a;
  color: #e4e4e4;
}

/* Container Styling */
.container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: #2c2c2c;
  padding: 30px;
  border-radius: 15px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.6);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.container:hover {
  transform: translateY(-5px) scale(1.01);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.8);
}

.small_logo {
  width: 110px;
  height: 110px;
  margin-bottom: 10px;
  background-image: url('../assets/images/mori.png');
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
  border-radius: 90%;
}

/* Title Styling */
.title {
  font-size: 30px;
  font-family: 'TITLE';
  font-weight: bold;
  margin-bottom: 20px;
  color: #e2e2e2;
  text-align: center;
  user-select: none;
}

/* Form Styling */
form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

form input {
  padding: 12px;
  width: 270px;
  border: none;
  border-radius: 8px;
  background: #444;
  color: #e4e4e4;
  font-size: 14px;
}

form input:focus {
  outline: 2px solid #9146bc;
}

/* Button Styling */
form button {
  padding: 12px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  background: #9146bc;
  color: #fff;
  font-size: 16px;
  transition: all 0.3s ease;
}

form button:hover {
  background: #7d3aa6;
}

/* Alternate Button Styling */
form button[type='button'] {
  background: transparent;
  color: #e4e4e4;
  border: 1px solid #444;
  transition: all 0.3s ease;
}

form button[type='button']:hover {
  background: #444;
}

.error-message {
  color: rgb(255, 77, 77);
  font-size: 14px;
  margin-top: 10px;
  text-align: center;
}

/* Popup Styling */
.popup {
  position: fixed;
  top: 15%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #9146bc;
  color: #e4e4e4;
  padding: 20px 40px;
  border-radius: 10px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.6);
  text-align: center;
  font-size: 16.5px;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.5s ease, transform 0.5s ease;
  animation: notification 1.5s ease;
}

.popup.hidden {
  opacity: 0;
  pointer-events: none;
  transform: translate(-50%, -60%);
}

.popup.show {
  opacity: 1;
  pointer-events: all;
  transform: translate(-50%, -50%);
}

.info_form {
  font-size: 15px;
  color: #e4e4e4;
  text-align: center;
  user-select: none;
}

#forgot_password {
  font-size: 13.5px;
  color: #e4e4e4;
  text-decoration: none;
  text-align: center;
  transition: all 0.2s;
  user-select: none;
}

#forgot_password:hover {
  cursor: pointer;
  color: #a75fd1;
}

@keyframes notification {
  0% {
    top: 5%;
  }
  100% {
    top: 15%;
  }
}
</style>
