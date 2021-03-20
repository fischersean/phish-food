<template>
  <nav class="flex flex-wrap items-center justify-between px-2 py-5 m-3">
    <div
      class="container px-4 mx-1 sm:mx-8 flex flex-wrap items-center justify-between"
    >
      <div>
        <router-link to="/">
          <img src="../assets/logo.svg" class="w-20" title="Home" />
        </router-link>
      </div>

      <div class="sm:flex flex-grow items-center">
        <ul class="flex flex-row list-none ml-5">
          <slot> </slot>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: "Nav",
};
</script>
