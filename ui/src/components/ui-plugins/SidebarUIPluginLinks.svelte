<script>
    import { link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
    import tooltip from "@/actions/tooltip";
    import { superuser } from "@/stores/superuser";
    import ApiClient from "@/utils/ApiClient";

    let plugins = [];

    $: if ($superuser?.id) {
        loadPlugins();
    }

    async function loadPlugins() {
        try {
            const loaded = await ApiClient.send("/api/ui-plugins", {
                method: "GET",
                requestKey: "ui-plugins",
            });
            plugins = loaded.plugins || []
        } catch (err) {
            if (!err?.isAbort) {
                console.warn("Failed to load UI plugins.", err);
            }
        }
    }
</script>

{#each plugins as plugin}
    <a
        href={`/ui-plugins/${plugin.base}/`}
        class="menu-item"
        aria-label={plugin.name}
        use:link
        use:active={{ path: `/ui-plugins/${plugin.base}/?.*`, className: "current-route" }}
        use:tooltip={{ text: plugin.name, position: "right" }}
    >
        <i class={plugin.icon} />
    </a>
{/each}
