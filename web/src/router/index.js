import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Hero",
    component: Home,
  },
  //{
  //path: "/docs",
  //name: "Documentation",
  //component: () => import("../views/Docs.vue"),
  //},
  //{
  //path: "/buy",
  //name: "Buy",
  //component: () => import("../views/Buy.vue"),
  //},
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

export default router;
