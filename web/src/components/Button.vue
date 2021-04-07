<template>
  <div>
    <button
      :onClick="onClick"
      class="rounded-full bg-green-200 p-3 focus:outline-white"
    >
      <slot> </slot>
    </button>
  </div>
</template>

<script>
export default {
  name: "Button",
  props: {
    onClick: String,
  },
};
</script>
