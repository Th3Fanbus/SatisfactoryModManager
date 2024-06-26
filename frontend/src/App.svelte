<script lang="ts">
  import './_global.postcss';
  import { arrow, autoUpdate, computePosition, flip, offset, shift, size } from '@floating-ui/dom';
  import { Modal, initializeStores, storePopup } from '@skeletonlabs/skeleton';
  import { setContextClient } from '@urql/svelte';

  import TitleBar from '$lib/components/TitleBar.svelte';
  import LeftBar from '$lib/components/left-bar/LeftBar.svelte';
  import ModDetails from '$lib/components/mod-details/ModDetails.svelte';
  import ErrorModal from '$lib/components/modals/ErrorModal.svelte';
  import ExternalInstallMod from '$lib/components/modals/ExternalInstallMod.svelte';
  import { supportedProgressTypes } from '$lib/components/modals/ProgressModal.svelte';
  import { modalRegistry } from '$lib/components/modals/modalsRegistry';
  import ImportProfile from '$lib/components/modals/profiles/ImportProfile.svelte';
  import ModsList from '$lib/components/mods-list/ModsList.svelte';
  import { initializeGraphQLClient } from '$lib/core/graphql';
  import { getModalStore, initializeModalStore } from '$lib/skeletonExtensions';
  import { installs, invalidInstalls, progress } from '$lib/store/ficsitCLIStore';
  import { error, expandedMod, siteURL } from '$lib/store/generalStore';
  import { konami } from '$lib/store/settingsStore';
  import { ExpandMod, GenerateDebugInfo, UnexpandMod } from '$wailsjs/go/app/app';
  import { Environment, EventsOn } from '$wailsjs/runtime';

  initializeStores();
  initializeModalStore();

  storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow, size });

  let frameless = false;
  Environment().then((env) => {
    if (env.buildType !== 'dev') {
      document.addEventListener('contextmenu', (event) => event.preventDefault());
    }
    if (env.platform === 'windows') {
      frameless = true;
    }
  });

  export let apiEndpointURL!: string;
  export let siteEndpointURL!: string;
  
  $: $siteURL = siteEndpointURL;

  setContextClient(initializeGraphQLClient(apiEndpointURL));

  let windowExpanded = false;

  $: if ($expandedMod) {
    ExpandMod().then(() => {
      setTimeout(() => {
        windowExpanded = true;
      }, 100);
    });
  } else {
    windowExpanded = false;
    setTimeout(() => {
      UnexpandMod();
    }, 100);
  }

  $: pendingExpand = $expandedMod && !windowExpanded;

  let invalidInstallsError = false;
  let noInstallsError = false;
  let focusOnEntry: HTMLSpanElement;

  const installsInit = installs.isInit;
  const invalidInstallsInit = invalidInstalls.isInit;

  $: if($installsInit && $invalidInstallsInit && $installs.length === 0) {
    if($invalidInstalls.length > 0) {
      invalidInstallsError = true;
    } else {
      noInstallsError = true;
    }
  } else {
    invalidInstallsError = false;
    noInstallsError = false;
  }

  const modalStore = getModalStore();

  $: if($progress && supportedProgressTypes.includes($progress.action)) {
    modalStore.triggerUnique({
      type: 'component',
      component: 'progress',
      meta: {
        persistent: true,
      },
    }, true);
  }

  $: if($error) {
    modalStore.trigger({
      type: 'component',
      component: {
        ref: ErrorModal,
        props: {
          error: $error,
        },
      },
    }, true);
    $error = null;
  }

  EventsOn('externalInstallMod', (modReference: string, version: string) => {
    if (!modReference) return;
    modalStore.trigger({
      type: 'component',
      component: {
        ref: ExternalInstallMod,
        props: {
          modReference,
          version,
        },
      },
    });
  });

  EventsOn('externalImportProfile', async (path: string) => {
    if (!path) return;
    modalStore.trigger({
      type: 'component',
      component: {
        ref: ImportProfile,
        props: {
          filepath: path,
        },
      },
    });
  });

  $: isPersistentModal = $modalStore.length > 0 && $modalStore[0].meta?.persistent;

  function modalMouseDown(event: MouseEvent) {
    if (!isPersistentModal) return;
    if (!(event.target instanceof Element)) return;
    const classList = event.target.classList;
    if (classList.contains('modal-backdrop') || classList.contains('modal-transition')) {
      event.stopImmediatePropagation();
    }
  }

  function modalKeyDown(event: KeyboardEvent) {
    if (!isPersistentModal) return;
    if (event.key === 'Escape') {
      event.stopImmediatePropagation();
    }
  }
  
  const code = [38, 38, 40, 40, 37, 39, 37, 39, 66, 65];
  const keyQueue: number[] = [];
  window.addEventListener('keydown', (event) => {
    keyQueue.push(event.keyCode);
    if (keyQueue.length > code.length) {
      keyQueue.shift();
    }
    if (keyQueue.length === code.length && keyQueue.every((val, idx) => code[idx] === val)) {
      $konami = !$konami;
    }
  });
</script>

<div class="flex flex-col h-screen w-screen select-none">
  {#if frameless}
    <TitleBar />
  {/if}
  <div class="flex grow h-0">
    <LeftBar />
    <div class="flex flex-auto @container/mod-list-wrapper z-[1]">
      <div class="{$expandedMod && !pendingExpand ? 'w-2/5 hidden @3xl/mod-list-wrapper:block @3xl/mod-list-wrapper:flex-auto' : 'w-full'}" class:max-w-[42.5rem]={!!$expandedMod}>
        <ModsList
          hideMods={noInstallsError || invalidInstallsError}
          on:expandedMod={() => {
            focusOnEntry.focus();
          }}>
          <div class="card my-auto mr-4">
            <header class="card-header font-bold text-2xl text-center">
              {#if noInstallsError}
                No Satisfactory installs found
              {:else}
                {$invalidInstalls.length} invalid Satisfactory install{$invalidInstalls.length !== 1 ? 's' : ''} found
              {/if}
            </header>
            <section class="p-4">
              <p class="text-base text-center">
                Seems wrong? Click the button below and send the generated zip file on the <a class="text-primary-600 underline" href="https://discord.gg/xkVJ73E">modding discord</a> in #help-using-mods.
              </p>
            </section>
            <footer class="card-footer">
              <button
                class="btn text-primary-600 w-full"
                on:click={GenerateDebugInfo}>
                Generate debug info
              </button>
            </footer>
          </div>
        </ModsList>
      </div>
      <div class="w-3/5" class:grow={!pendingExpand} class:hidden={!$expandedMod}>
        <ModDetails bind:focusOnEntry/>
      </div>
    </div>
  </div>
</div>

<!--
  skeleton modals don't provide a way to make them persistent (i.e. ignore mouse clicks outside and escape key)
  but we can capture the events and stop them if the modal has the persistent meta flag set, and the event would have closed the modal
-->
<svelte:window
  on:keydown|capture|nonpassive={modalKeyDown}
  on:mousedown|capture|nonpassive={modalMouseDown} />
<Modal components={modalRegistry} />
