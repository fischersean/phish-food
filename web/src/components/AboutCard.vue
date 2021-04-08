<template>
  <div>
    <div class="flex flex-col text-center">
      <h1 class="text-2xl font-bold mb-5">
        Trending data covering all of the major stock-trading subreddits
      </h1>
      <p class="mb-2 text-xl">
        Access to the same information available to the more "sophisticated"
        investors, at the same time
      </p>
      <Button
        onClick="parent.open('https://rapidapi.com/fischersean/api/thekettle/')"
        class="mt-5"
        >Get API Access →</Button
      >
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
