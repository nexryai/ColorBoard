<script lang="ts">
    import defaultThumbnail from "$lib/images/default-thumbnail.avif"
    import { Button } from "$lib/components/ui/button"
    import { Skeleton } from "$lib/components/ui/skeleton"
    import BrandGoogle from "@tabler/icons-svelte/icons/brand-google"
    import BrandAzure from "@tabler/icons-svelte/icons/brand-azure"
    import AddGalleryButton from "$lib/components/AddGalleryButton.svelte"
    import { getApps, initializeApp, type FirebaseApp } from "firebase/app"
    import {
        getAuth,
        getIdToken,
        GoogleAuthProvider,
        signInWithPopup,
        signInWithRedirect,
        onAuthStateChanged,
        setPersistence,
        inMemoryPersistence,
        type Persistence,
    } from "firebase/auth"

    import { isLoggedIn } from "$lib/account"
    import { type Gallery, type Image, fetchMyGalleries } from "$lib/api"
    import { getFirebaseConfig } from "$lib/firebase"
    import { goto } from "$app/navigation"

    let isLoading = true
    let galleries: Gallery[]

    function getImageSrc(images?: Image[]) {
        if (!images || images.length === 0) {
            return defaultThumbnail
        } else {
            return `/api/files/${images[0].thumbnailKey}`
        }
    }

    if (isLoggedIn()) {
        console.log("Fetching galleries...")
        fetchMyGalleries()
            .then((res) => {
                galleries = res
                isLoading = false
            })
            .catch((error) => {
                console.error(error)
            })
    }

    let app: FirebaseApp | undefined

    async function checkToken() {
        const firebaseConfig = await getFirebaseConfig()

        if (!getApps().length) {
            app = initializeApp(firebaseConfig)
        }

        const auth = getAuth(app)

        onAuthStateChanged(auth, (user) => {
            if (user) {
                // ログインしていれば中通る
                console.log(user) // ユーザー情報が表示される
                getIdToken(user, false).then(token => {console.log(token)})
            } else {
                console.log("not logged in")
            }
        })

        /*
        const user = auth.currentUser;
        if (!user) {
            console.log("user is null");
        } else {
            const token = await getIdToken(user, false);
            console.log(token);
        }*/
    }

    async function loginWithGoogle() {
        try {
            const firebaseConfig = await getFirebaseConfig()

            if (!getApps().length) {
                app = initializeApp(firebaseConfig)
            }

            const auth = getAuth(app)
            const provider = new GoogleAuthProvider()
            // デフォルトだとリフレッシュトークンがindexedDBに保存される
            // ブラウザ側で平文でディスクに保存されるのは避けたいのでリフレッシュトークンは使わないようにする（JWTの期限が切れたら再認証）
            // 連携ログインを使ってるので2回目以降の再認証はユーザーの操作がなくても数回のリダイレクトでできる
            setPersistence(auth, inMemoryPersistence).then(() => {
                signInWithPopup(auth, provider)
            })
        } catch (e) {
            console.log(e)
        }
    }

    checkToken()
</script>

<section>
    {#if !isLoggedIn()}
        <div class="login-form">
            <div class="box">
                <p class="box-title">Welcome back</p>
                <div class="box-bottom">
                    <div class="form-large-icon">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 24 24"
                            fill="currentColor"
                            ><path
                                d="M1 5C1 4.44772 1.44772 4 2 4H22C22.5523 4 23 4.44772 23 5V19C23 19.5523 22.5523 20 22 20H2C1.44772 20 1 19.5523 1 19V5ZM13 8V10H19V8H13ZM18 12H13V14H18V12ZM10.5 10C10.5 8.61929 9.38071 7.5 8 7.5C6.61929 7.5 5.5 8.61929 5.5 10C5.5 11.3807 6.61929 12.5 8 12.5C9.38071 12.5 10.5 11.3807 10.5 10ZM8 13.5C6.067 13.5 4.5 15.067 4.5 17H11.5C11.5 15.067 9.933 13.5 8 13.5Z"
                            ></path></svg
                        >
                    </div>
                    <div class="login-links">
                        <div class="login-button">
                            <Button
                                class="w-[280px]"
                                on:click={async () => await loginWithGoogle()}
                            >
                                <BrandGoogle />
                                　Login with Google
                            </Button>
                        </div>
                        <div class="input-elm">
                            <Button
                                class="w-[280px]"
                                on:click={async () => await loginWithGoogle()}
                            >
                                <BrandAzure />
                                　Login with Microsoft Entra ID
                            </Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {:else}
        <div class="center-container">
            <div class="flex justify-between items-center">
                <p class="text-2xl">Galleries</p>
                <AddGalleryButton />
            </div>

            {#if !isLoading}
                <div
                    class="grid place-items-center gap-4 grid-cols-2 md:grid-cols-4 lg:grid-cols-5"
                >
                    {#each galleries as gallery}
                        <a href="/galleries/{gallery.id}" class="w-[150px]">
                            <div
                                class="cursor-pointer h-[150px] w-[150px] mt-[24px] overflow-hidden shadow-md rounded-2xl transition hover:shadow-lg"
                            >
                                <img
                                    class="gallery-box h-[150px] w-[150px] aspect-auto object-cover"
                                    src={getImageSrc(gallery.images)}
                                    alt=""
                                />
                            </div>
                            <p class="w-[150px] mt-2 truncate">
                                {gallery.name}
                            </p>
                        </a>
                    {/each}
                </div>
            {:else}
                <Skeleton class="h-[150px] w-[250px] mt-8" />
            {/if}
        </div>
    {/if}
</section>

<style lang="scss">
    .login-form {
        width: 720px;
        border: solid #e1e1e1 1px;
        border-radius: 10px;
        margin: 100px auto 50px auto;
        overflow: hidden;

        /*boxに背景画像を設定*/
        .box {
            position: relative;

            background-image: url("https://images.unsplash.com/photo-1655635643617-72e0b62b9278?q=80&w=2832&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D");
            background-size: cover;
            background-position: center;
            height: 370px;
            display: flex;
            justify-content: center;
            align-items: center;
            flex-direction: column;

            .box-title {
                font-size: 24px;
                color: white;
                font-weight: bold;
                position: absolute;
                top: 0;
                left: 0;
                margin: 5px 0 0 32px;
            }

            .box-bottom {
                background-color: #fff;
                padding: 20px;
                width: 100%;
                height: calc(100% - 50px);
                /*底に配置*/
                position: absolute;
                bottom: 0;

                .form-large-icon {
                    font-size: 64px;
                    color: #606060;
                    margin: 22px auto;

                    svg {
                        margin: auto;
                        width: 45px;
                        height: 45px;
                    }
                }

                .login-links {
                    width: 370px;
                    margin: auto;
                    text-align: center;

                    .login-button {
                        margin: 8px;
                    }
                }
            }
        }
    }

    .center-container {
        text-align: left;
        min-height: 85vh;
    }
</style>
