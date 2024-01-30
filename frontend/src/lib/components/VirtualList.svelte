<script generics="T" lang="ts">
  import { tick } from 'svelte';

  // eslint-disable-next-line no-undef
  export let items: T[];
  export let itemHeight: number | undefined = undefined;
  export let bench = 10;
  let clazz = '';
  export { clazz as class };
  export let containerClass = '';

  let start = 0;
  let end = 10;

  let viewport: HTMLElement;
  let container: HTMLElement;
  let heightMap: number[] = [];
  
  $: {
    heightMap = Array.from({ length: items.length });
    if(viewport && container && items.length > 0) {
      tick().then(updateVisible);
    }
  }


  function updateHeightMap() {
    const virtualRows = Array.from(container?.children ?? []);
    virtualRows.forEach((elem, idx) => {
      heightMap[start + idx] = elem.clientHeight;
    });
    
    const setHeights = heightMap.filter((x) => !!x);
    const averageHeight = setHeights.reduce((acc, curr) => acc + curr, 0) / setHeights.length;
    
    heightMap = heightMap.map((x) => x ?? (itemHeight ?? averageHeight));
  }

  let viewportHeight: number;

  $: if(viewport) {
    // Add or remove elements when the viewport height changes
    viewportHeight;
    updateVisible();
  }

  async function updateVisible() {
    const { scrollTop } = viewport;

    updateHeightMap();
    
    let height = 0;
    let newStart = 0;
    while(newStart < items.length && height + heightMap[newStart] < scrollTop) {
      height += heightMap[newStart];
      newStart++;
    }

    let newEnd = newStart;
    while(newEnd < items.length && height < scrollTop + viewport.clientHeight) {
      height += heightMap[newStart];
      newEnd++;
    }

    newStart = Math.max(newStart - bench, 0);
    newEnd = Math.min(newEnd + bench, items.length);

    const old_start = start;

    start = newStart;
    end = newEnd;
    
    // prevent jumping if we scrolled up
    // this became an issue when using using skeleton popups as tooltips,
    // specifically when calling window.getComputedStyle(tooltipElement).<anyPropertyHere>
    // I have no idea why that causes this issue
    if (start < old_start) {
      await tick();
      viewport.scrollTo({ top: scrollTop, behavior: 'instant' });
    }
  }

  $: top = heightMap.slice(0, start).reduce((acc, curr) => acc + curr, 0);
  $: bottom = heightMap.slice(end).reduce((acc, curr) => acc + curr, 0);
  $: visibleItems = items.map((data, index) => ({ index, data })).slice(start, end);
</script>

<div
  bind:this={viewport}
  class="relative overflow-y-scroll h-full {clazz}"
  bind:offsetHeight={viewportHeight}
  on:scroll={updateVisible}
>
  <div
    bind:this={container}
    style="padding-top: {top}px; padding-bottom: {bottom}px;"
    class={containerClass}
  >
    {#each visibleItems as item (item.index)}
      <div class="overflow-hidden">
        <slot item={item.data}>Missing template</slot>
      </div>
    {/each}
  </div>
</div>