<template>
  <div class="wrapper">
    <router-link class="logout-link" v-if="loggedIn" to="/logout">Log out</router-link>

    <NewApartment :showModal="this.showModal"></NewApartment>

    <div class="dashboard">
      <div class='sidebar'>
        <div class='heading'>
          <h1>Apartments</h1>
          <v-btn slot="activator" class="new-apartment" fab dark small color="blue"
                 v-if="canCrudApartment()" @click="showModal = true;">
            <v-icon dark>add</v-icon>
          </v-btn>
        </div>
        <div class="filters">
          <form class="filters-form" @submit.prevent="filterApartments()">
            <v-text-field label="Area (m2)" type="number" step="0.001"
                          v-model.number="filterData.floorAreaMeters"></v-text-field>
            <v-text-field label="Price (USD)" type="number" min="0.0" step="0.01"
                          v-model.number="filterData.pricePerMonthUSD"></v-text-field>
            <v-text-field label="Rooms" type="number" min="0"
                          v-model.number="filterData.roomCount"></v-text-field>
            <v-btn fab dark small type="submit" color="black">
              <v-icon dark>search</v-icon>
            </v-btn>
          </form>
        </div>
        <div id='listings' class='listings'>
          <v-card class="item" v-for="rental of rentals" :key="rental.id">
            <a href="#" @click="panTo(rental.latitude, rental.longitude)">
              {{ rental.name }}
            </a><em v-if="!rental.available">(Occupied)</em>
            <div class="detail">
              <b>Price:</b> <em> ${{ rental.pricePerMonthUSD }} </em>
              <b>Area:</b> <em> {{ rental.floorAreaMeters }}m2 </em>
              <b>Rooms:</b> <em> {{ rental.roomCount }} </em>
              <div class="added"><b>Added:</b> <em> {{ formatDate(rental.dateAdded) }}</em></div>
              <div class="desc" v-if="rental.description">{{ rental.description }}</div>
            </div>
            <v-btn outline @click="toggleAvailability(rental)"
                   v-if="canCrudApartment()">{{ rental.available ? "Rent out" : "Available"}}
            </v-btn>
          </v-card>
        </div>
      </div>

      <GmapMap class="map" ref="mmm" :center="mapProps.center" :zoom="mapProps.zoom">
        <GmapInfoWindow :options="mapProps.infoWindowOptions" :position="infoWindowPos" :opened="infoWindowOpen"
                        @closeclick="infoWindowOpen = false">
          {{ infoContent }}
        </GmapInfoWindow>


        <GmapMarker :key="i" v-for="(m, i) in markers" :position="m.position"
                    @click="toggleInfoWindow(m, i)"></GmapMarker>
      </GmapMap>
    </div>
  </div>
</template>


<script>
    import $auth from "./auth";
    import $rentals from './rentals';
    import $users from "./users";
    import NewApartment from './NewApartment';

    export default {
        name: 'Dashboard',
        components: {NewApartment},
        data() {
            return {
                mapProps: {
                    center: {
                        lat: 0,
                        lng: 0,
                    },
                    zoom: 2,
                    infoWindowOptions: {
                        pixelOffset: {
                            width: 0,
                            height: -35
                        }
                    }
                },
                rentals: [],
                infoWindowPos: null,
                infoWindowOpen: false,
                infoContent: "",


                filterData: {
                    floorAreaMeters: null,
                    pricePerMonthUSD: null,
                    roomCount: null,
                },

                userData: {
                    username: null,
                    role: null,
                },

                allUsers: [],

                showModal: false
            }
        },

        created() {
            this.loadApartments();
            this.loadUserData();
            this.tryLoadingAllUsers();
        },

        computed: {
            loggedIn() {
                return $auth.isLoggedIn();
            },

            markers() {
                const m = [];
                for (let s of this.rentals) {
                    m.push({
                        position: {
                            lat: s.latitude,
                            lng: s.longitude
                        },
                        infoText: `name: ${s.name}\n price: $ ${s.pricePerMonthUSD}`,
                    });
                }
                return m;
            },

            // Return a list of users formated for v-select
            allAdminsAndRealtors() {
                const realtors = [];
                for (let u of this.allUsers) {
                    if (u.role === "realtor" || u.role === "admin") {
                        realtors.push({
                            text: u.username,
                            value: u.id
                        });
                    }
                }

                return realtors;
            }
        },

        methods: {
            toggleInfoWindow(marker, idx) {
                this.infoWindowPos = marker.position;
                this.infoContent = marker.infoText;
                if (this.currentMidx === idx) {
                    this.infoWindowOpen = !this.infoWindowOpen;
                } else {
                    this.infoWindowOpen = true;
                    this.currentMidx = idx;
                }
            },

            panTo(lat, lng) {
                this.$refs.mmm.panTo({
                    lat: lat,
                    lng: lng
                });
                this.mapProps.zoom = 5;
            },

            panOut() {
                this.mapProps.center = {lat: 0, lng: 0};
                this.mapProps.zoom = 2;
            },


            filterApartments() {
                $rentals.loadAllApartments(this.filterData).then(res => {
                    this.rentals = res;
                    this.panOut();
                });
            },

            loadApartments() {
                $rentals.loadAllApartments({}).then(res => {
                    this.rentals = res;
                });
            },

            loadUserData() {
                $auth.getUserInfo().then(res => {
                    this.userData = res.data;
                }).catch(err => {
                    alert(err);
                })
            },

            tryLoadingAllUsers() {
                $users.getAllUsers().then(res => {
                    this.allUsers = res;
                }).catch(err => {
                    console.log(err);
                })
            },

            canCrudApartment() {
                return this.userData.role === "admin" || this.userData.role === 'realtor';
            },

            toggleAvailability(apartment) {
                $rentals.changeAvailability(apartment.id, !apartment.available).then(() => {
                    this.loadApartments();
                }).catch(err => {
                    alert(err)
                });
            },

            formatDate(strDate) {
                const d = new Date(strDate);
                return d.toLocaleDateString("en-us", {day: "numeric", month: "long", year: "numeric"});
            }
        },
    }
</script>

<style scoped>
  .dashboard {
    display: grid;
    grid-template-columns: 30% 70%;
    height: 100vh;
  }

  .sidebar {
    border-right: 1px solid rgba(0, 0, 0, 0.25);
    overflow: hidden;
    height: 100vh;
  }


  h1 {
    font-size: 22px;
    margin: 0;
    font-weight: 400;
    line-height: 20px;
    padding: 20px 2px;
  }

  .heading {
    display: grid;
    border-bottom: 1px solid #eee;
    grid-template-columns: 1fr 1fr;
    min-height: 60px;
    line-height: 60px;
    padding: 0 10px;
  }

  .sidebar .new-apartment {
    justify-self: end;
    text-align: center;
    position: relative;
    margin: .8rem .5rem;
    right: 0;
  }

  .new-apartment-form {
    display: grid;
    grid-template-rows: repeat(auto-fit, 1fr);
    row-gap: .3rem;
    padding: .5rem;
  }

  .new-apartment-form > input, .new-apartment-form > button {
    height: 2rem;
    font-size: .9rem;
    width: 100%;
  }

  .new-apartment-form > textarea {
    height: 4rem;
    font-size: .9rem;
    width: 100%;
  }

  .filters-form {
    display: grid;
    grid-template-columns: repeat(3, 4fr) 1fr;
  }

  .filters-form > input {
    width: 100%;
  }

  .logout-link {
    position: absolute;
    top: 18px;
    left: 170px;
  }

  .listings {
    height: 90%;
    overflow: auto;
  }

  .listings .item {
    display: block;
    border-bottom: 1px solid #eee;
    padding: 10px;
    text-decoration: none;
  }

  .detail .desc {
    color: #555;
  }
</style>
