<template>
  <div>
    <div class="text-xl block px-4 py-2">
      <router-link :to="href.to" exact-active-class="font-bold underline">
        {{ href.label }}
        <slot> </slot>
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: "NavItem",
  props: {
    href: {
      label: String,
      to: String,
    },
  },
  computed: {
    isActive: function () {
      return this.$route.path == this.href.to;
    },
  },
};
</script>
