function showPopup() {
    const popup = document.getElementById("popup");
  
    popup.classList.remove("hidden");
    popup.classList.add("show");
  
    // Hide the popup after 3 seconds
    setTimeout(() => {
      popup.classList.remove("show");
      popup.classList.add("hidden");
    }, 3000);
  }
  
  document
    .getElementById("reset-password-form")
    .addEventListener("submit", async (e) => {
      e.preventDefault();
      const errorDiv = document.getElementById("reset-password-error");
      errorDiv.style.display = "none";
  
      // Get the values from the inputs
      const token = new URLSearchParams(window.location.search).get("token");
      const newPassword = document.getElementById("new-password").value;
      const confirmPassword = document.getElementById("confirm-password").value;
  
      // Validate if passwords match
      if (newPassword !== confirmPassword) {
        errorDiv.textContent = "Passwords do not match.";
        errorDiv.style.display = "block";
        return;
      }
  
      // Send the request to the server
      const response = await fetch("/verify-reset-token", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token, new_password: newPassword }),
      });
  
      const result = await response.json();
      if (!response.ok) {
        errorDiv.textContent = result.error || "An error occurred.";
        errorDiv.style.display = "block";
      } else {
        const new_password = document.getElementById("new-password");
        const confirm_password = document.getElementById("confirm-password");
        new_password.value = "";
        confirm_password.value = "";
        showPopup();
        setTimeout(() => {
          window.location.href = "/";
        }, 2000);
      }
    });
  