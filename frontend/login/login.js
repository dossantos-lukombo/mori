// JavaScript to toggle between login and register forms
function toggleForms() {
  const loginForm = document.getElementById("login-form");
  const registerForm = document.getElementById("register-form");

  if (loginForm.style.display === "none") {
    loginForm.style.display = "flex";
    registerForm.style.display = "none";
  } else {
    loginForm.style.display = "none";
    registerForm.style.display = "flex";
  }
}

function backToLogin() {
  const loginForm = document.getElementById("login-form");
  const reset_password_form = document.getElementById("reset-password-form");

  loginForm.style.display = "flex";
  reset_password_form.style.display = "none";
}

function resetPasswordButton() {
  const loginForm = document.getElementById("login-form");
  const reset_password_form = document.getElementById("reset-password-form");

  loginForm.style.display = "none";
  reset_password_form.style.display = "flex";
}

// Utility to get a cookie by name
function getCookie(name) {
  const cookies = document.cookie.split("; ");
  for (let i = 0; i < cookies.length; i++) {
    const [key, value] = cookies[i].split("=");
    if (key === name) return value;
  }
  return null;
}

// Function to handle login form submission
document.getElementById("login-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const errorDiv = document.getElementById("login-error");
  errorDiv.style.display = "none";

  // Get the CSRF token from the cookie
  const csrfToken = getCookie("csrf_token"); // Ensure this is the right cookie

  if (!csrfToken) {
    errorDiv.textContent = "CSRF token is missing.";
    errorDiv.style.display = "block";
    return;
  }

  const formData = new FormData(e.target);
  formData.append("csrf_token", csrfToken);
  const response = await fetch("/login", {
    method: "POST",
    body: formData,
  });

  const result = await response.text();
  if (!response.ok) {
    errorDiv.textContent = result.error || "An error occurred during login.";
    errorDiv.style.display = "block";
  } else {
    console.log("Login successful");
    console.log("Redirecting to home chat...");
    console.log("result : ", result);
    console.log("response : ", response);
    console.log("response.body : ", response.body);
    console.log("url path : ", window.location.pathname);
    window.location.href = response.url; // Redirect to dashboard on successful login


  }
});

function showPopup(message) {
  const popup = document.getElementById("popup");
  const popupMessage = document.getElementById("popup-message");

  popupMessage.textContent = message;
  popup.classList.remove("hidden");
  popup.classList.add("show");

  // Hide the popup after 3 seconds
  setTimeout(() => {
    popup.classList.remove("show");
    popup.classList.add("hidden");
  }, 3000);
}

// Register form submission
document
  .getElementById("register-form")
  .addEventListener("submit", async (e) => {
    e.preventDefault();
    const errorDiv = document.getElementById("register-error");
    errorDiv.style.display = "none";

    // Get the CSRF token from the cookie
    const csrfToken = getCookie("csrf_token"); // Ensure this is the right cookie

    if (!csrfToken) {
      errorDiv.textContent = "CSRF token is missing.";
      errorDiv.style.display = "block";
      return;
    }

    const formData = new FormData(e.target);
    formData.append("csrf_token", csrfToken); // Append the CSRF token to the form data

    try {
      const response = await fetch("/register", {
        method: "POST",
        body: formData,
      });

      const result = await response.json(); // Always parse the response as JSON
      if (!response.ok) {
        errorDiv.textContent =
          result.error || "An error occurred during registration.";
        errorDiv.style.display = "block";
        if (result.reloadCaptcha) {
          reloadCaptcha();
        }
      } else {
        showPopup("Registration successful, please verify your email.");
        toggleForms(); // Switch to login form
      }
    } catch (err) {
      errorDiv.textContent = err.message;
      errorDiv.style.display = "block";
    }
  });

// Function to reload captcha
function reloadCaptcha() {
  const captchaImage = document.getElementById("captcha-image");
  captchaImage.src = "/captcha?" + new Date().getTime(); // Append timestamp to prevent caching
}

document
  .getElementById("reset-password-form")
  .addEventListener("submit", async (e) => {
    e.preventDefault();
    const errorDiv = document.getElementById("reset-password-error");
    errorDiv.style.display = "none";

    const formData = new FormData(e.target);
    const response = await fetch("/reset-password", {
      method: "POST",
      body: formData,
    });

    const result = await response.json();
    if (!response.ok) {
      errorDiv.textContent =
        result.error || "An error occurred during password reset.";
      errorDiv.style.display = "block";
    } else {
      showPopup("Password reset link sent to your Email.");
      backToLogin();
    }
  });