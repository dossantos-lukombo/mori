import { createRouter, createWebHistory } from "vue-router";
import Auth from "../components/Auth.vue";
import store from "@/store";

const routes = [
  {
    path: "/",
    name: "auth",
    component: Auth,
  },
  {
    path: "/sign-in",
    name: "sign-in",
    component: () => import("../views/SignInView.vue"),
  },
  {
    path: "/register",
    name: "register",
    component: () => import("../views/RegisterView.vue"),
  },
  {
    path: "/main",
    name: "mainpage",
    components: {
      default: () => import("../views/MainView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue"),
    },
  },
  {
    path: "/messages",
    name: "messages",
    component: () => import("../components/ChatBoxView.vue"),
    props: (route) => ({
      name: route.query.name || "Conversation",
      receiverId: route.query.receiverId || null,
      type: route.query.type || "PERSON",
    }),
  },
  {
    path: "/profile/:id",
    name: "Profile",
    components: {
      default: () => import("../views/ProfileView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue"),
    },
  },
  {
    path: "/group/:id",
    name: "Group",
    components: {
      default: () => import("../views/GroupView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue"),
    },
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach(async (to, from) => {
  const isAuthenticated = await store.dispatch("isLoggedIn");

  // Redirect to sign-in if not authenticated, except for sign-in and register
  if (!isAuthenticated && to.name !== "sign-in" && to.name !== "register") {
    return { name: "sign-in" };
  }
});

export default router;
