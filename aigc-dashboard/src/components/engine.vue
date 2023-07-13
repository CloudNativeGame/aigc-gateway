<script>

const NOT_INSTALLED = "Not Installed"
const RUNNING = "Running"
const PAUSED = "Paused"
const PENDING = "Pending"

export default {
  props: ["engine", "userinfo"],
  data: function () {
    return {
      dialog: false,
      state: NOT_INSTALLED,
      server: {},
    }
  },
  methods: {
    getData: function () {
      let engine = this.engine;
      let name = engine.metadata.name;
      let namespace = engine.metadata.namespace;

      let customData = this.userinfo.custom_data;

      let key = namespace + "-" + name;

      if (customData == undefined) {
        this.state = NOT_INSTALLED
      }

      if (customData[key] != undefined) {

        let meta = customData[key];
        let id = meta["ID"];

        this.axios.get("/resource/" + namespace + "/" + name + "/" + id).then((response) => {
          this.server = response.data

          let status = this.server["status"]["currentState"]
          let networkStatus = this.server["status"]["networkStatus"]["currentNetworkState"]

          if (status == "Ready" && networkStatus == "Ready") {
            this.state = RUNNING
          } else {
            this.state = PENDING
          }

        }).catch((error) => {
          if (error.response.status == 423) {
            this.state = PAUSED
          }
        })
      }
    },
    visit: function () {
      if (this.server != undefined) {
        let networkStatus = this.server["status"]["networkStatus"]
        let networkType = networkStatus["networkType"]
        let externalAddresses = networkStatus["externalAddresses"]
        if (networkType == "Kubernetes-Ingress") {
          window.open("http://" + externalAddresses[0].endPoint, "'_blank'")
        } else {
          window.open("http://" + externalAddresses[0].ip + ":" + externalAddresses[0]["ports"][0]["port"], "'_blank'")
        }
      }
    },
    confirm: function () {
      this.dialog = true;
    },
    create: function () {
      let engine = this.engine;
      let name = engine.metadata.name;
      let namespace = engine.metadata.namespace;

      this.axios.put("/resource/" + namespace + "/" + name).then((response) => {
        this.items = response.data
        this.dialog = false;
      }).catch((error) => {
        this.dialog = false;
      })
    },
    pause: function () {
      let engine = this.engine;
      let name = engine.metadata.name;
      let namespace = engine.metadata.namespace;

      this.axios.post("/resource/" + namespace + "/" + name + "/pause").then((response) => {
        this.items = response.data
      }).catch((error) => {
      })
    },
    recover: function () {
      let engine = this.engine;
      let name = engine.metadata.name;
      let namespace = engine.metadata.namespace;

      this.axios.post("/resource/" + namespace + "/" + name + "/recover").then((response) => {
        this.items = response.data
      }).catch((error) => {
      })
    },
    del: function () {
      let engine = this.engine;
      let name = engine.metadata.name;
      let namespace = engine.metadata.namespace;

      this.axios.delete("/resource/" + namespace + "/" + name).then((response) => {
          this.items = response.data
      }).catch((error) => {
      })
    },
    restart: function () {
      let engine = this.engine;
      let name = engine.metadata.name;
      let namespace = engine.metadata.namespace;

      this.axios.post("/resource/" + namespace + "/" + name + "/restart").then((response) => {
          this.items = response.data
      }).catch((error) => {
      })
    }
  },
  created: function () {
    this.getData();
  },
  mounted: function () {
    return

    // let customData = this.userinfo.custom_data;
    //
    // if (customData == undefined) {
    //   return
    // }
    //
    // for (let key in customData) {
    //   let namespacedName = key.split("-")
    //   this.axios.post("/resource/" + namespacedName[0] + "/" + namespacedName[1] + "/recover").then((response) => {
    //     console.log("recover " + key);
    //   }).catch((error) => {
    //     console.log(error)
    //   })
    // }
  }
}
</script>

<template>
  <v-col cols="4">

    <v-card variant="outlined">
      <img style="width:30%;display:block;margin: 0 auto;margin-top:20px;" src="@/assets/logo.svg"/>
      <v-container>
        <v-card-text>
          <div style="text-align:center;font-weight: bolder;font-size: 24px">{{ engine.metadata.name }}</div>
        </v-card-text>
      </v-container>

      <v-container>
        <v-row class="text-center" v-if="state==='Pending'">
          <v-col>
            <v-progress-circular
                indeterminate
                color="primary"
                text-align="center"
                style="height:48px"
            ></v-progress-circular>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="6" v-if="state==='Not Installed'">
            <v-btn flexbox height="48" width="100%" color="indigo-darken-3"
                   @click="confirm()">
              Install
            </v-btn>
          </v-col>
          <v-col cols="6" v-if="state==='Paused'">
            <v-btn flexbox height="48" width="100%" color="indigo-darken-3" @click="recover()">
              Recover
            </v-btn>
          </v-col>

          <v-col cols="6" v-if="state==='Running'">
            <v-btn flexbox height="48" width="100%" color="indigo-darken-3" @click="visit()">
              Visit
            </v-btn>
          </v-col>
          <v-col cols="6" v-if="state==='Running'">
            <v-btn flexbox height="48" width="100%" color="indigo-darken-3" @click="pause()">
              Pause
            </v-btn>
          </v-col>
          <v-col cols="6" v-if="state!=='Not Installed'">
            <v-btn flexbox height="48" width="100%" color="indigo-darken-3" @click="del()">
              Uninstall
            </v-btn>
          </v-col>
          <v-col cols="6" v-if="state!=='Not Installed' && state!=='Paused'">
            <v-btn flexbox height="48" width="100%" color="indigo-darken-3" @click="restart()">
              Restart
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>
  </v-col>

  <div class="text-center">
    <v-dialog
        v-model="dialog"
        width="auto">
      <v-card>
        <v-card-text>
          Installation will create new resource,would you like to do that.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" @click="create()">OK</v-btn>
          <v-btn color="primary" @click="dialog = false">CLOSE</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>