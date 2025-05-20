/**
 * Helper function to get the correct avatar URL based on the path
 * Handles both old imageUpload paths and new imageUploads paths
 * Also ensures default avatars are handled correctly
 * 
 * @param {string} avatarPath - The path to the avatar image
 * @returns {string} - The correct URL to display the avatar
 */
export function getAvatarUrl(avatarPath) {
  if (!avatarPath) {
    return 'http://localhost:8081/frontend/public/images/default.svg';
  }

  // Check if the path already contains the server URL
  if (avatarPath.startsWith('http')) {
    return avatarPath;
  }

  // Handle the default image path format
  if (avatarPath === 'frontend/public/images/default.svg') {
    return 'http://localhost:8081/frontend/public/images/default.svg';
  }

  // Handle the legacy path (imageUpload without 's')
  if (avatarPath.startsWith('imageUpload/')) {
    // Convert to new path format
    return `http://localhost:8081/imageUploads/${avatarPath.substring(12)}`;
  }

  // Handle the new path (imageUploads with 's')
  if (avatarPath.startsWith('imageUploads/')) {
    return `http://localhost:8081/${avatarPath}`;
  }

  // Default case: assume it's a relative path
  return `http://localhost:8081/${avatarPath}`;
} 