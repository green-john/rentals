<template>
    <div class="wrapper">
        <router-link class="logout-link" v-if="loggedIn" to="/logout">Log out</router-link>
        <div class="dashboard">
            <div class='sidebar'>
                <div class='heading'>
                    <h1>Apartments</h1>
                    <button v-if="canCrudApartment()" @click="showModal()" class="new-apartment">
                        +
                    </button>
                </div>
                <div class="filters">
                    <form class="filters-form" @submit.prevent="filterApartments()">
                        <input placeholder="Area (m2)" type="number" step="0.001"
                               v-model.number="filterData.floorAreaMeters">
                        <input placeholder="Price (USD)" type="number" min="0.0" step="0.01"
                               v-model.number="filterData.pricePerMonthUSD">
                        <input placeholder="Rooms" type="number" min="0"
                               v-model.number="filterData.roomCount">
                        <input type="submit" value="Y">
                    </form>
                </div>
                <div id='listings' class='listings'>
                    <div class="item" v-for="rental of rentals" :key="rental.id">
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
                        <button @click="toggleAvailability(rental)"
                                v-if="canCrudApartment()">{{ rental.available ? "Rent out" : "Available"}}
                        </button>
                    </div>
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

            <modal name="new-apartment" @before-close="clearApartmentData" width="300px"
                   height="auto" pivotY.number="0.2">
                <form @submit.prevent="createApartment()" class="new-apartment-form">
                    <h2 style="text-align: center">Create Apartment</h2>
                    <input placeholder="Name" required type="text" v-model.trim="newApartmentData.name">
                    <textarea placeholder="Description (optional)"
                              v-model.trim="newApartmentData.description"></textarea>
                    <input placeholder="Floor Area (m2)" required type="number" step="0.001"
                           v-model.number="newApartmentData.floorAreaMeters">
                    <input placeholder="Price per month (USD)" required type="number" min="0.0" step="0.01"
                           v-model.number="newApartmentData.pricePerMonthUSD">
                    <input placeholder="Room count" type="number" required min="0"
                           v-model.number="newApartmentData.roomCount">
                    <input placeholder="Latitude" type="number" required step="0.0000001" min="-90" max="90"
                           v-model.number="newApartmentData.latitude">
                    <input placeholder="Longitude" type="number" required step="0.0000001" min="-180" max="180"
                           v-model.number="newApartmentData.longitude">
                    <input placeholder="Realtor Id" type="number" required min="0"
                           v-model.number="newApartmentData.realtorId">
                    <label for="available">Available
                        <input id="available" type="checkbox" v-model="newApartmentData.available">
                    </label>

                    <input type="submit" value="Create">
                    {{ newApartmentMessage }}
                </form>
            </modal>
        </div>
    </div>
</template>


<script>
    import $auth from "./auth";
    import $rentals from './rentals';

    export default {
        name: 'Dashboard',
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

                newApartmentData: {
                    name: null,
                    description: null,
                    realtorId: null,
                    floorAreaMeters: null,
                    pricePerMonthUSD: null,
                    roomCount: null,
                    latitude: null,
                    longitude: null,
                    available: null,
                },

                filterData: {
                    floorAreaMeters: null,
                    pricePerMonthUSD: null,
                    roomCount: null,
                },

                userData: {
                    username: null,
                    role: null,
                },

                newApartmentMessage: ""
            }
        },

        created() {
            this.loadApartments();
            this.loadUserData();
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

            showModal() {
                this.$modal.show("new-apartment");
            },

            createApartment() {
                $rentals.newApartment(this.newApartmentData).then(res => {
                    alert(`Apartment created with id ${res.id}`);
                    this.clearApartmentData();
                    this.loadApartments();
                }).catch(err => {
                    this.newApartmentMessage = `Error: ${err}`;
                });
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

            canCrudApartment() {
                return this.userData.role === 'admin' || this.userData.role === 'realtor';
            },

            toggleAvailability(apartment) {
                $rentals.changeAvailability(apartment.id, !apartment.available).then(() => {
                    this.loadApartments();
                }).catch(err => {
                    alert(err)
                });
            },

            clearApartmentData() {
                this.newApartmentData = {
                    name: null,
                    description: null,
                    realtorId: null,
                    floorAreaMeters: null,
                    pricePerMonthUSD: null,
                    roomCount: null,
                    latitude: null,
                    longitude: null,
                    available: null,
                };
                this.newApartmentMessage = "";
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
        background-color: white;
        border: none;
        font-size: 1.5rem;
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
