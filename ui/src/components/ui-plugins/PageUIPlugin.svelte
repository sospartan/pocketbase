<script>
    
    import { location } from "svelte-spa-router";
    import { pageTitle } from "@/stores/app";

    $pageTitle = "plugin page";

</script>

<div class="page-wrapper">
    {#if $location}
    <iframe
        src={$location}
        title="Plugin Manager"
        class="plugin-iframe"
    />
    {/if}
</div>

<style lang="scss">
    .plugin-iframe {
        width: 100%;
        height: 100%;
        border: none;
    }
</style>


