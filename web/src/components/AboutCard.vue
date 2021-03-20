<template>
  <div class="py-5">
    <div class="flex flex-col">
      <h1 class="text-2xl font-bold mb-5">
        Trending data covering all of the major stock-trading subreddits
      </h1>
      <p class="mb-2 text-xl">
        Access to the same information available to the more "sophisticated"
        investors
      </p>
      <Button class="py-10">Get API Access →</Button>
    </div>
  </div>
</template>

<script>
import Button from "@/components/Button.vue";
export default {
  name: "AboutCard",
  components: {
    Button,
  },
};
</script>
