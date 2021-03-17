<template>
  <div class="font-black text-6xl flex flex-row">
    <img class="mr-2" src="../assets/logo.svg" width="40" />
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
export default {
  name: "Hero",
  props: {
    msg: String,
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
