<script setup lang="ts">
	import { ref, onMounted, defineProps, defineEmits } from 'vue';
	import { ThreadStats, ThreadNavProps } from "@/model/thread/thread.model.ts";
	
	const autoTimer = ref<int>(10);
	const isAutoTimerUsed = ref<bool>(false);
	const autoTimerCountdown = () => {
		setTimeout(() => {
			if (autoTimer.value <= 0) {
				autoTimer.value = 10;
				emit('threadUpdate');
			} else {
				if (isAutoTimerUsed.value) {
					autoTimer.value -= 1;
				}
			}
			
			autoTimerCountdown();
		}, 1000);	
	}

	const props = defineProps<ThreadNavProps>();
	const emit = defineEmits(['threadUpdate']);
	
	const doUpdate = () => {
		emit('threadUpdate');
	}
	
	onMounted(() => {
		autoTimerCountdown();
	});
</script>

<template>
<hr />
<nav>
	<span class="left">
		[<RouterLink class="link" :to="`/${props.board_code}`">Return</RouterLink>]
		[<RouterLink class="link" :to="`/${props.board_code}/catalog`">Catalog</RouterLink>]
		<template v-if="props.jump_to_id && props.jump_to_label">
			[<a class="link" :href="`#${props.jump_to_id}`">{{ props.jump_to_label }}</a>]
		</template>
		[<a class="link" href="#" @click.prevent="doUpdate">Update</a>]
		[<input type="checkbox" v-model="isAutoTimerUsed" name="auto"><label for="auto"> Auto</label>] <template v-if="isAutoTimerUsed">{{ autoTimer }}</template>
	</span>
	<span class ="right">
		<span class="tooltip">{{ props.thread_stats.posts }}
			<span class="tooltiptext">Replies</span>
		</span> / 
		<span class="tooltip">{{ props.thread_stats.images }}
			<span class="tooltiptext">Images</span>
		</span> / 
		<span class="tooltip">{{ props.thread_stats.posters }}
			<span class="tooltiptext">Users</span>
		</span> / 
		<span class="tooltip">{{ props.thread_stats.page }}
			<span class="tooltiptext">Page</span>
		</span>
	</span>
</nav>
<hr />
</template>

<style scoped>
	RouterLink, a {
		color: #34345c;
	}
	
	RouterLink:hover, a:hover {
		color: #DD0000;
	}
	
	nav {
		display: flex;
		justify-content: space-between;
		
		.right {
			cursor: default;
		}
	}
</style>