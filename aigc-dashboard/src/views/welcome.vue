<script setup>
</script>

<script>
export default {
  data: () => ({
    alignments: [
      'start',
      'center',
      'end',
    ],
  }),
  methods: {
    signIn: function () {
      window.location.href = "/sign-in";
      return false
    }
  },
}
</script>

<template>
  <v-container style="padding-top:15%">
    <v-row no-gutters>
      <v-col cols="8">
        <v-sheet style="display: flex;justify-content: center; align-items: center; height: 100%;">
          <div style="height:250px;">
            <div style="font-size: 72px;font-weight: bold;text-align: center">AIGC Gateway</div>
            <div style="font-size: 32px;font-weight: bold;text-align: center">OpenSource AIGC engine on kubernetes.
            </div>
            <v-btn @click="signIn" block height="62" color="indigo-darken-3" min-width="164" style="margin:20px auto">
              Start to Use
            </v-btn>
          </div>
        </v-sheet>
      </v-col>
      <v-col>
        <v-sheet>
          <v-img
              src="https://img.alicdn.com/imgextra/i4/O1CN010LpTnx1JoqzbZl7sV_!!6000000001076-2-tps-1272-1072.png"></v-img>
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>
