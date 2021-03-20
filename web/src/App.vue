<template>
  <div id="app" class="font-mono flex flex-col">
    <!--<Nav>-->
    <!--<NavItem :href="{ to: '/docs', label: 'Docs' }"> </NavItem>-->
    <!--<NavItem :href="{ to: '/buy', label: 'Buy' }"> </NavItem>-->
    <!--</Nav>-->
    <router-view class="self-center" />
  </div>
</template>

<script>
//import Nav from "./components/Nav.vue";
//import NavItem from "./components/Navitem.vue";

export default {
  name: "App",
  components: {
    //Nav,
    //NavItem,
  },
  data: function () {
    return {};
  },
};
</script>

<style></style>
