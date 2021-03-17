<template>
  <div id="app" class="font-mono flex flex-col">
    <div class="self-center items-center">
      <div class="mt-10 mb-10 text-xl" id="nav">
        <router-link to="/">Home</router-link> |
        <router-link to="/about">About</router-link> |
        <router-link to="/docs">Documentation</router-link>
      </div>
    </div>
    <router-view class="self-center" />
  </div>
</template>

<style></style>
