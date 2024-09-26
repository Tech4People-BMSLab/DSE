// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------


// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
export const sleep = (ms) => {
    return new Promise(resolve => setTimeout(resolve, ms))
}

export const wait_until = async (predicate, timeout=-1, interval=100) => {
    const time_start = Date.now()
    while (true) {
        const result = await predicate()
        if (result) {
            break
        }
        if (timeout > 0 && Date.now() - time_start > timeout) {
            throw new Error('TimeoutError')
        }
        await sleep(interval)
    }
}
