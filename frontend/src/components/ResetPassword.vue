<template>
    <div class="container">
      <div class="small_logo"></div>
      <div class="title">Reset Password</div>
      <form @submit.prevent="resetPassword">
        <p class="info_form">Enter your new password.</p>
        <input type="password" v-model="newPassword" placeholder="New Password" required />
        <input type="password" v-model="confirmPassword" placeholder="Confirm Password" required />
        <button type="submit">Reset Password</button>
        <div v-if="resetError" class="error-message">{{ resetError }}</div>
      </form>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        newPassword: "",
        confirmPassword: "",
        resetError: "",
      };
    },
    methods: {
      async resetPassword() {
        if (this.newPassword !== this.confirmPassword) {
          this.resetError = "Passwords do not match.";
          return;
        }
  
        try {
          const response = await fetch("/verify-reset-token", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ token: this.$route.query.token, new_password: this.newPassword }),
          });
          if (!response.ok) {
            const result = await response.json();
            this.resetError = result.error || "An error occurred.";
          } else {
            this.$router.push("/");
          }
        } catch (error) {
          this.resetError = "An error occurred.";
        }
      },
    },
  };
  </script>
  
  <style scoped>
  @font-face {
  font-family: "TITLE";
  src: url('@/assets/fonts/SinosansRegular-aYxZ5.otf') format("truetype");
  font-weight: normal;
  font-style: normal;
}

/* General Body Styling */
body {
  font-family: "Arial", sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  margin: 0;
  background-color: #1a1a1a; /* Dark background */
  color: #e4e4e4; /* Light text for contrast */
}

/* Container Styling */
.container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: #2c2c2c; /* Slightly lighter dark background */
  padding: 30px;
  border-radius: 15px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.6);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.container:hover {
  transform: translateY(-5px) scale(1.01); /* Lift effect and slight scale */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.8);
}

/* Small Logo Styling */
.small_logo {
  width: 110px;
  height: 110px;
  margin-bottom: 10px;
  background-image: url('@/assets/images/mori.png');
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
  border-radius: 50%;
}

/* Title Styling */
.title {
  font-size: 30px;
  font-family: "TITLE";
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
  background: #444; /* Input background */
  color: #e4e4e4; /* Text inside input */
  font-size: 14px;
}

form input:focus {
  outline: 2px solid #9146bc; /* Highlight color */
}

/* Button Styling */
form button {
  padding: 12px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  background: #9146bc; /* Primary button color */
  color: #fff; /* Button text color */
  font-size: 16px;
  transition: all 0.3s ease;
}

form button:hover {
  background: #7d3aa6; /* Hover effect for primary button */
}

/* Error Message Styling */
#reset-password-error {
  color: rgb(255, 77, 77);
  font-size: 14px;
  margin-top: 10px;
  text-align: center;
  display: none; /* Initially hidden */
}


.info_form{
  font-size: 15px;
  color: #e4e4e4;
  text-align: center;
  user-select: none;
}

/* Popup Styling */
.popup {
  position: fixed;
  top: 15%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #9146bc; /* Slightly lighter dark background */
  color: #e4e4e4;
  padding: 20px 40px;
  border-radius: 10px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.6);
  text-align: center;
  font-size: 16.5px;
  z-index: 1000; /* Ensure it appears above other elements */
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

@keyframes notification {
  0% {
    top: 5%;
  }
  100% {
    top: 15%;
  }
  
}
</style>