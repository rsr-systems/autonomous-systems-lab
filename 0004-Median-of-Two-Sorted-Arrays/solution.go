package main

import "math"

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    if len(nums1) > len(nums2) {
        return findMedianSortedArrays(nums2, nums1)
    }

    m, n := len(nums1), len(nums2)
    low, high := 0, m

    for low <= high {
        partitionX := (low + high) / 2
        partitionY := (m+n+1)/2 - partitionX

        maxLeftX := math.MinInt32
        if partitionX > 0 {
            maxLeftX = nums1[partitionX-1]
        }

        minRightX := math.MaxInt32
        if partitionX < m {
            minRightX = nums1[partitionX]
        }

        maxLeftY := math.MinInt32
        if partitionY > 0 {
            maxLeftY = nums2[partitionY-1]
        }

        minRightY := math.MaxInt32
        if partitionY < n {
            minRightY = nums2[partitionY]
        }

        if maxLeftX <= minRightY && maxLeftY <= minRightX {
            if (m+n)%2 == 0 {
                leftMax := math.Max(float64(maxLeftX), float64(maxLeftY))
                rightMin := math.Min(float64(minRightX), float64(minRightY))
                return (leftMax + rightMin) / 2.0
            }
            return math.Max(float64(maxLeftX), float64(maxLeftY))
        } else if maxLeftX > minRightY {
            high = partitionX - 1
        } else {
            low = partitionX + 1
        }
    }

    return 0.0
}
