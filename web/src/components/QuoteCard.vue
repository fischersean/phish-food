<template>
  <div class="py-5 font-serif text-xl">
    <p class="italic">
      <slot></slot>
    </p>
    <p class="font-bold mt-2">- {{ attribution }}</p>
  </div>
</template>

<script>
export default {
  name: "QuoteCard",
  props: ["attribution"],
};
</script>
