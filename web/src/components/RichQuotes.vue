<template>
  <div>
    <QuoteCard attribution="Leon Cooperman, Hedge Fund Manager">
      "Fair share is a bullshit concept. It is just a way of attacking wealthy
      people"
    </QuoteCard>
    <QuoteCard attribution="Jim Cramer, Theranos Supporter">
      "Never underestimate the power of enthusiastic buyers who do not know what
      they’re doing"
    </QuoteCard>
  </div>
</template>

<script>
import QuoteCard from "@/components/QuoteCard.vue";

export default {
  name: "RichQuotes",
  components: {
    QuoteCard,
  },
};
</script>
