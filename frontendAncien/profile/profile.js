document.addEventListener("DOMContentLoaded", function () {
    // Change Avatar
    document.getElementById("avatar-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const response = await fetch("/profile/avatar", {
            method: "POST",
            body: formData,
        });
        if (!response.ok) {
            console.error("Error updating avatar");
        } else {
            alert("Avatar updated successfully!");
        }
    });

    // Update User Info
    document.getElementById("update-info-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const response = await fetch("/profile/update", {
            method: "PUT",
            body: formData,
        });
        if (!response.ok) {
            console.error("Error updating information");
        } else {
            alert("User information updated successfully!");
        }
    });

    // Change Password
    document.getElementById("password-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const response = await fetch("/profile/change-password", {
            method: "PUT",
            body: formData,
        });
        if (!response.ok) {
            console.error("Error changing password");
        } else {
            alert("Password changed successfully!");
        }
    });

    // Delete Account
    document.getElementById("delete-account-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const response = await fetch("/profile/delete-account", {
            method: "DELETE",
        });
        if (!response.ok) {
            console.error("Error deleting account");
        } else {
            alert("Account deleted successfully!");
            window.location.href = "/";
        }
    });

    // Delete All Favorites
    document.getElementById("delete-favorites-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const response = await fetch("/profile/delete-favorites", {
            method: "DELETE",
        });
        if (!response.ok) {
            console.error("Error deleting favorites");
        } else {
            alert("Favorites deleted successfully!");
        }
    });

    // Delete All Conversations
    document.getElementById("delete-conversations-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const response = await fetch("/profile/delete-conversations", {
            method: "DELETE",
        });
        if (!response.ok) {
            console.error("Error deleting conversations");
        } else {
            alert("Conversations deleted successfully!");
        }
    });
});
