/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package helper

func GenCopyright(comment ...string) string {
	copyrightComment := `/* 
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */
`

	if len(comment) > 0 {
		copyrightComment += "\n"
		for _, c := range comment {
			copyrightComment += c + "\n"
		}
	}
	copyrightComment += "\n"

	return copyrightComment
}
