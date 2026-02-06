import { defineStore } from 'pinia';
import { GetAlbumList } from '../../wailsjs/go/core/AppCore';

const DEFAULT_ALBUM = { id: 0, name: '默认相册' };

export const useBaseStore = defineStore('base', {
	state: () => ({
		albumList: [DEFAULT_ALBUM],
		hasLoadedAlbumList: false,
	}),
	actions: {
		/**
		 * 拉取相册列表（带缓存与幂等保护）
		 * - 避免多次调用导致重复追加
		 * - 优先读取 sessionStorage 缓存
		 * - 缓存缺失时调用后端接口
		 */
		async fetchAlbumList() {
			// 幂等保护：已加载则直接返回
			if (this.hasLoadedAlbumList) {
				return;
			}
			const cacheKey = 'albumListCache';
			try {
				// 优先读取缓存，避免不必要的接口请求
				const cached = sessionStorage.getItem(cacheKey);
				if (cached) {
					const parsed = JSON.parse(cached);
					if (Array.isArray(parsed)) {
						// 使用缓存数据覆盖列表，并保留默认相册
						this.albumList = [
							DEFAULT_ALBUM,
							...parsed.map((item: any) => ({
								id: item.id,
								name: item.name,
							})),
						];
						this.hasLoadedAlbumList = true;
						return;
					}
				}

				// 缓存不存在或无效时请求后端
				const albums = await GetAlbumList();
				const { status, data } = albums;
				if (status === true && Array.isArray(data)) {
					// 使用接口数据覆盖列表，并写入缓存
					this.albumList = [
						DEFAULT_ALBUM,
						...data.map((item: any) => ({
							id: item.id,
							name: item.name,
						})),
					];
					sessionStorage.setItem(cacheKey, JSON.stringify(data));
					this.hasLoadedAlbumList = true;
				}
			} catch (error) {
				console.error('Failed to fetch album list:', error);
			}
		},
        // 格式化时间
        formatDateTimeToSecond(value: string){
            const text = (value || '').trim()
            if (!text) return '-'

            // 已经是 "YYYY-MM-DD HH:mm:ss" 则直接返回
            if (!text.includes('T') && text.includes(' ')) return text

            const d = new Date(text)
            if (Number.isNaN(d.getTime())) return text

            const pad2 = (n: number) => String(n).padStart(2, '0')
            const yyyy = d.getFullYear()
            const mm = pad2(d.getMonth() + 1)
            const dd = pad2(d.getDate())
            const hh = pad2(d.getHours())
            const mi = pad2(d.getMinutes())
            const ss = pad2(d.getSeconds())
            return `${yyyy}-${mm}-${dd} ${hh}:${mi}:${ss}`
        }
	},
});