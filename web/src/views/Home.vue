<template>
  <div class="sm:p-16 w-3/4 flex flex-col self-center">
    <Hero />
    <div class="flex flex-col sm:py-10">
      <AboutCard class="py-10" />
      <ExampleResponse />
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import Hero from "@/components/Hero.vue";
import AboutCard from "@/components/AboutCard.vue";
import ExampleResponse from "@/components/ExampleResponse.vue";

export default {
  name: "Home",
  components: {
    Hero,
    AboutCard,
    ExampleResponse,
  },
};
</script>
