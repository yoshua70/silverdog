<script lang="ts">
  import JobStatus from "$lib/jobStatus";

  export let jobTitle: string;
  export let jobStatus: string;
  export let jobErrorDetails = "";
  export let jobBody = "";

  const getBackgroundColor = () => {
    switch (jobStatus.toLowerCase()) {
      case "pending":
        return "bg-gray-300";
      case "inprogress":
        return "bg-yellow-300";
      case "completed":
        return "bg-green-300";
      case "failed":
        return "bg-red-300";
      default:
        return "bg-gray-300";
    }
  };
</script>

<div class={`p-8 flex flex-col gap-2 ${getBackgroundColor()} rounded-xl`}>
  <p>Name: {jobTitle}</p>
  <p>Status: <span class="uppercase">{jobStatus}</span></p>
  <p>Body: {jobBody}</p>
  {#if jobStatus === JobStatus.Failed}
    <p>Details: {jobErrorDetails}</p>
  {/if}
</div>
