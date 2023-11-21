# Algorithms Supported

`smash` supports a variety of hashing algorithms.

<table>
<thead>
 <tr>
    <th>Algorithm</th>
    <th>Default</th>
    <th>Variations / Aliases</th>
</tr>
</thead>
<tbody>
    <tr>
        <td>
            xxhash<br/>
            <sub><sup><a href="https://xxhash.com/">learn more</a></sup></sub>
        </td>
        <td><code>xxhash</code></td>
        <td>
            <ul>
                <li><code>xxhash</code></li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>
            murmur3<br/>
            <sub><sup><a href="https://en.wikipedia.org/wiki/MurmurHash">learn more</a></sup></sub>
        </td>
        <td><code>murmur3</code></td>
        <td>
            <ul>
                <li><code>murmur3</code> (alias: <code>murmur3-128</code>)</li>
                <li><code>murmur3-128</code></li>
                <li><code>murmur3-64</code></li>
                <li><code>murmur3-32</code></li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>SHA-256</td>
        <td><code>sha256</code></td>
        <td>
            <ul>
                <li><code>sha256</code></li>
                <li><code>sha-256</code></li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>SHA-512</td>
        <td><code>sha512</code></td>
        <td>
            <ul>
                <li><code>sha512</code></li>
                <li><code>sha-512</code></li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>MD5</td>
        <td><code>md5</code></td>
        <td>
            <ul>
                <li><code>md5</code></li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>FNV128<br/>
        <sub><sup><a href="https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function">learn more</a></sup></sub></td>
        <td><code>fnv128</code></td>
        <td>
            <ul>
                <li><code>fnv128</code></li>
                <li><code>fnv-128</code></li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>FNV128a<br/>
            <sub><sup><a href="https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function">learn more</a></sup></sub></td>
        <td><code>fnv128a</code></td>
        <td>
            <ul>
                <li><code>fnv</code> (alias: <code>fnv128a</code>)</li>
                <li><code>fnv128a</code></li>
                <li><code>fnv-128a</code></li>
            </ul>
        </td>
    </tr>
</tbody>
</table>

Generally, when slicing is enabled (default), we'd recommend `xxhash` or `murmur3`. 

When you're wanting a full hash (`--disable-slicing` option), generally `sha512` or `sha-256`.
