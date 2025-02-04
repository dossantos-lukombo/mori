<template>
  <div>
    <!-- NavBar on top -->
    <NavBarOn />
    <!-- UpgradeIA component (your drop file area) -->
    <UpgradeIA />
  </div>
</template>

<script>
import NavBarOn from "@/components/NavBarOn.vue";
import UpgradeIA from "@/components/UpgradeIA.vue";

export default {
  name: "DropFiles",
  components: {
    NavBarOn,
    UpgradeIA,
  },
};
</script>

<style scoped>
/* Optional: Add view-specific styling if needed */
</style>
