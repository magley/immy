<script setup lang="ts">
	import { ref, onMounted, defineProps, defineEmits } from 'vue';
	
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
	
	interface Props {
		board_code: string,
		jump_to_id: string | null,
		jump_to_label: string | null,
	}
	
	const props = defineProps<Props>();
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
	[<RouterLink class="link" :to="`/${props.board_code}`">Return</RouterLink>]
	[<RouterLink class="link" :to="`/${props.board_code}/catalog`">Catalog</RouterLink>]
	<template v-if="props.jump_to_id && props.jump_to_label">
		[<a class="link" :href="`#${props.jump_to_id}`">{{ props.jump_to_label }}</a>]
	</template>
	[<a class="link" href="#" @click.prevent="doUpdate">Update</a>]
	[<input type="checkbox" v-model="isAutoTimerUsed"> Auto] <template v-if="isAutoTimerUsed">{{ autoTimer }}</template>
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
</style>