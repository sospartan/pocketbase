<script>
    
    import { location } from "svelte-spa-router";
    import { pageTitle } from "@/stores/app";
    import {
        collections,
        loadCollections
    } from "@/stores/collections";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";

    $pageTitle = "plugin page";

    let collectionUpsertPanel;
    let iframe;

    function listenIframe() {
        window.addEventListener("message", (event) => {
            
            const { data } = event;
            const { source, action} = data;
            if(source.startsWith("ui-plugin/")) {
                if (action === "edit-collection") {
                    const { id } = data;
                    showCollectionUpsertPanel(id);
                }
            }
            
        });
    }

    async function showCollectionUpsertPanel(id) {
        if (!$collections) {
            await loadCollections();
        }

        const collection = $collections.find((c) => c.id == id);
        if (collection) {
            collectionUpsertPanel.show(collection);
        }
    }

    function notifyChanged() {
        loadCollections();
        iframe.contentWindow.postMessage({
            source: "ui-plugin-host",
            action: "changed",
        }, "*");
    }

    
    listenIframe();
    loadCollections();

</script>

<div class="page-wrapper">
    {#if $location}
    <iframe
        bind:this={iframe}
        src={$location}
        title="Plugin Manager"
        class="plugin-iframe"
    />
    {/if}
</div>
<CollectionUpsertPanel
    bind:this={collectionUpsertPanel}
    on:delete={() => {
        notifyChanged()
    }}
    on:save={() => {
        notifyChanged()
    }}
/>

<style lang="scss">
    .plugin-iframe {
        width: 100%;
        height: 100%;
        border: none;
    }
</style>


