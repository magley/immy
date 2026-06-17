<script setup lang="ts">
	import { BanAPI, type CreateBanDTO } from '@/api/ban.api';
	import type { BoardDTO } from '@/api/board.api';
	import type { PostDTO } from '@/api/post.api';
	import { GetPostTimeReadable } from '@/model/post/post.model';
	import { DateFromDuration } from '@/util/various.util';
	import { AxiosError } from 'axios';
	import { onMounted } from 'vue';
	import { ref, useId } from 'vue';

	interface CreateBanProps {
		post: PostDTO | undefined;
		currentBoard: BoardDTO | undefined;
	}

	const props = defineProps<CreateBanProps>();
	const error = ref<string | undefined>(undefined);
	const userID = ref<number | undefined>(undefined);


	const id = useId();
	const banDTO = ref<CreateBanDTO>({
		ip_start: '',
		ip_end: null,
		expires_at: null,
		board_id: null,
		reason: '',
		warning: false
	});
	const formPermanentBan = ref<boolean>(false);
	const banExpirationPresetVal = ref<string>("5min");
	const banDurationVal = ref<string>("");
	const isRangeban = ref<boolean>(false);

	onMounted(() => {
		applyBanDuration();
	});

	const onClickSubmitBan = () => {
		submitBan();
	}

	const onClickSetBanDuration = () => {
		applyBanDuration();
	}

	const submitBan = () => {
		if (banDTO.value.warning) {
			// Cleaning up, although it doesn't matter. If the ban is a
			// warning, it must ignore these fields.x
			banDTO.value.expires_at = null;
			banDTO.value.board_id = null;
		}

		error.value = undefined;

		console.log(banDTO.value.ip_start);

		BanAPI.CreateBan(banDTO.value).then((res) => {

		}).catch((err: AxiosError) => {
			error.value = "Could not ban user";
			console.log(err);
		});
	}

	const applyBanDuration = () => {
		let durationValueFinal = banExpirationPresetVal.value;
		if (banExpirationPresetVal.value == "custom") {
			durationValueFinal = banDurationVal.value
		}
		banDTO.value.expires_at = DateFromDuration(durationValueFinal).toISOString();
	}
</script>

<template>
	<div>
		<div>
			<label :for="`warning-${id}`">Warning:</label>
			<input type=checkbox :id="`warning-${id}`" v-model="banDTO.warning" />
		</div>

		<!-- Ban only -->
		<div v-if="!banDTO.warning">
			<!-- Ban IP -->
			<div v-if="!props.post">
				<div>
					<label :for="`is-range-${id}`">Is range ban:</label>
					<input type=checkbox :id="`is-range-${id}`" v-model="isRangeban" />
				</div>

				<div v-if="!isRangeban">
					<label :for="`ip-${id}`">IP address:</label>
					<input :id="`ip-${id}`" v-model="banDTO.ip_start" placeholder="172.20.0.1" />
				</div>
				<div v-else>
					<label :for="`ip-start-${id}`">IP range start:</label>
					<input :id="`ip-start-${id}`"  v-model="banDTO.ip_start" placeholder="172.20.0.1" />
					<br/>
					<label :for="`ip-end-${id}`">IP range end:</label>
					<input :id="`ip-end-${id}`"  v-model="banDTO.ip_end" placeholder="172.20.0.254" />
				</div>
			</div>

			<!-- Ban duration -->
			<div>
				<div>
					<label :for="`warning-${id}`">Permanent:</label>
					<input type=checkbox :id="`permanent-${id}`" v-model="formPermanentBan" />
				</div>

				<div v-if="formPermanentBan">
					<b>This ban is permanent.</b>
				</div>
				<div v-else>
					<div>
						Expires on: {{ GetPostTimeReadable(banDTO.expires_at!) }}
					</div>
					<div>
						<b>Set duration:</b>
						<select v-model="banExpirationPresetVal">
							<option value="5min">5 minutes</option>
							<option value="1h">1 hour</option>
							<option value="24h">24 hours</option>
							<option value="3d">3 days</option>
							<option value="7d">7 days</option>
							<option value="28d">28 days</option>
							<option value="1y">1 year</option>
							<option value="custom">Custom</option>
						</select>
						<br/>

						<div v-if="banExpirationPresetVal == 'custom'">
							<label :for="`ban-duration-custom-${id}`">Enter duration:</label>
							<input
							:id="`ban-duration-custom-${id}`"
							v-model="banDurationVal"
							placeholder="5min, 1h, 5 days, 1 year" />
						</div>

						<button @click="onClickSetBanDuration">Set ban duration</button>
					</div>

				</div>
			</div>
			<!-- Board -->
			<div v-if="props.currentBoard">
				<label :for="`ban-local-${id}`">Ban on /{{ props.currentBoard.code }}/</label>
				<input :id="`ban-local-${id}`" type="radio" v-model="banDTO.board_id" :value="props.currentBoard.id" />
				<br/>
				<label :for="`ban-global-${id}`">Ban globally</label>
				<input :id="`ban-global-${id}`" type="radio" v-model="banDTO.board_id" :value="null" />
			</div>
		</div>

		<!-- Warning only -->
		<div v-else>

		</div>

		<!-- Common -->
		<div>

			<textarea v-model="banDTO.reason" placeholder="Ban message" rows="4" cols="20"></textarea>

			<br/>
			<div class="error" v-if="error">{{ error }}</div>
			<button @click="onClickSubmitBan">Submit ban</button>
		</div>
	</div>
</template>

<style>
	.error {
		color: var(--user-error-color);
	}
</style>